package signing

import "time"

type SigningTime struct {
	v time.Time
}

func SigningTimeOf(t time.Time) SigningTime {
	return SigningTime{
		// AWS S3 specify the explicit use of UTC time
		v: t.UTC(),
	}
}

func (st SigningTime) ShortFormat() string {
	const formatYYYYMMDD = "20060102"
	return st.v.Format(formatYYYYMMDD)
}

func (st SigningTime) LongFormat() string {
	const formatXAmzDate = "20060102T150405Z"
	return st.v.Format(formatXAmzDate)
}
