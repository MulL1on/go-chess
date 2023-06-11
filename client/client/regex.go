package client

import "regexp"

var (
	RegexTypeMove      = regexp.MustCompile(`^[a-h][1-8][a-h][1-8]$`)
	RegexTypeCastle    = regexp.MustCompile(`^O-O(-O)?$`)
	RegexTypePromotion = regexp.MustCompile(`^[a-h][18][a-h][18]=[QRBN]$`)
)
