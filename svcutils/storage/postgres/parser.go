package postgres

func IsNilStr(in string) *string {
	if in == "" {
		return nil
	}

	return &in
}
