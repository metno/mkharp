package harp

type Parameter struct {
	Parameter  string
	AccumHours float32
	Units      string
}

func GetParameter(name string) *Parameter {
	switch name {
	case "T2m":
		return &Parameter{
			Parameter:  "T2m",
			AccumHours: 0,
			Units:      "degC",
		}
	case "Tmax":
		return &Parameter{
			Parameter:  "Tmax",
			AccumHours: 0,
			Units:      "degC",
		}
	case "Tmin":
		return &Parameter{
			Parameter:  "Tmin",
			AccumHours: 0,
			Units:      "degC",
		}
	case "AccPcp1h":
		return &Parameter{
			Parameter:  "AccPcp1h",
			AccumHours: 1,
			Units:      "kg/m^2",
		}

	case "AccPcp6h":
		return &Parameter{
			Parameter:  "AccPcp6h",
			AccumHours: 6,
			Units:      "kg/m^2",
		}

	case "AccPcp12h":
		return &Parameter{
			Parameter:  "AccPcp12h",
			AccumHours: 12,
			Units:      "kg/m^2",
		}

	case "AccPcp24h":
		return &Parameter{
			Parameter:  "AccPcp24h",
			AccumHours: 24,
			Units:      "kg/m^2",
		}
	}
	return nil
}
