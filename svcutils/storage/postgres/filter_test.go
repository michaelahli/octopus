package postgres

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildFilter(t *testing.T) {
	type args struct {
		query      string
		filter     map[string][]string
		exactQuery []SpecialQuery
	}

	type result struct {
		Query string
		Qargs []interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    result
		wantErr bool
	}{
		{
			name: "Positive Test Case: No Filter",
			args: args{
				query:  "SELECT * FROM transfer WHERE %s",
				filter: map[string][]string{},
			},
			want: result{
				Query: "SELECT * FROM transfer  ",
				Qargs: nil,
			},
			wantErr: false,
		},
		{
			name: "Positive Test Case: Normal Case",
			args: args{
				query: "SELECT * FROM transfer WHERE %s ORDER BY transfer.created DESC",
				filter: map[string][]string{
					"transfer.initiator_org": {"123456"},
					"transfer.rail":          {"DIRECT"},
					"transfer.status":        {"REJECTED", "ERROR", "REVERTED"},
				},
			},
			want: result{
				Query: "SELECT * FROM transfer WHERE transfer.initiator_org IN ($1) AND transfer.rail IN ($2) AND transfer.status IN ($3, $4, $5) ORDER BY transfer.created DESC",
				Qargs: []interface{}{
					"123456", "DIRECT", "REJECTED", "ERROR", "REVERTED",
				},
			},
			wantErr: false,
		},
		{
			name: "Positive Test Case: With Special Query",
			args: args{
				query: "select * from loan where %s",
				filter: map[string][]string{
					"loan_type": {"4"},
					"status":    {"INACTIVE"},
				},
				exactQuery: []SpecialQuery{
					{
						Query: "customer_id IN (select cbs_id from customer where branch_id = ?)",
						Args:  []string{"12345"},
					},
				},
			},
			want: result{
				Query: "select * from loan where loan_type IN ($1) AND status IN ($2) AND customer_id IN (select cbs_id from customer where branch_id = $3)",
				Qargs: []interface{}{
					"4", "INACTIVE", "12345",
				},
			},
			wantErr: false,
		},
		{
			name: "Positive Test Case: With Multiple Special Query",
			args: args{
				query: "select * from loan where %s",
				filter: map[string][]string{
					"loan_type": {"4"},
					"status":    {"INACTIVE"},
				},
				exactQuery: []SpecialQuery{
					{
						Query: "created > ?",
						Args:  []string{"2006-02-01"},
					},
					{
						Query: "created < ?",
						Args:  []string{"2006-02-02"},
					},
				},
			},
			want: result{
				Query: "select * from loan where loan_type IN ($1) AND status IN ($2) AND created > $3 AND created < $4",
				Qargs: []interface{}{
					"4", "INACTIVE", "2006-02-01", "2006-02-02",
				},
			},
			wantErr: false,
		},
		{
			name: "Positive Test Case: Only Special Query",
			args: args{
				query:  "select * from loan where %s",
				filter: map[string][]string{},
				exactQuery: []SpecialQuery{
					{
						Query: "created > ?",
						Args:  []string{"2006-02-01"},
					},
					{
						Query: "created < ?",
						Args:  []string{"2006-02-02"},
					},
				},
			},
			want: result{
				Query: "select * from loan where created > $1 AND created < $2",
				Qargs: []interface{}{
					"2006-02-01", "2006-02-02",
				},
			},
			wantErr: false,
		},
		{
			name: "Negative Test Case: Invalid Args",
			args: args{
				query: "select * from loan where %s",
				filter: map[string][]string{
					"loan_type": {"1"},
					"status":    {""},
				},
				exactQuery: []SpecialQuery{
					{
						Query: "created > ?",
						Args:  []string{""},
					},
				},
			},
			want: result{
				Query: "select * from loan where loan_type IN ($1)",
				Qargs: []interface{}{
					"1",
				},
			},
			wantErr: false,
		},
		{
			name: "Positive Test Case: Only Special Query With Wultiple Args",
			args: args{
				query:  "select * from loan where %s",
				filter: map[string][]string{},
				exactQuery: []SpecialQuery{
					{
						Query: "created > ? AND created < ?",
						Args:  []string{"2006-02-01", "2006-02-02"},
					},
				},
			},
			want: result{
				Query: "select * from loan where created > $1 AND created < $2",
				Qargs: []interface{}{
					"2006-02-01", "2006-02-02",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, qargs, err := BuildFilter(tt.args.query, tt.args.filter, tt.args.exactQuery...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildFilter err: %v wantErr: %v", err, tt.wantErr)
				return
			}

			got := result{
				Query: query,
				Qargs: qargs,
			}

			if !cmp.Equal(tt.want, got) {
				t.Errorf("BuildFilter diff: %s", cmp.Diff(tt.want, got))
				return
			}
		})
	}
}
