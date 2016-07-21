package hlp

// HerePlacesResponse represent results of search request to HERE API
type HerePlacesResponse struct {
	Results struct {
		Items []struct {
			AverageRating float64 `json:"averageRating,omitempty"`
			Category      struct {
				Href   string `json:"href,omitempty"`
				ID     string `json:"id,omitempty"`
				System string `json:"system,omitempty"`
				Title  string `json:"title,omitempty"`
				Type   string `json:"type,omitempty"`
			} `json:"category,omitempty"`
			Distance int       `json:"distance,omitempty"`
			Href     string    `json:"href,omitempty"`
			Icon     string    `json:"icon,omitempty"`
			ID       string    `json:"id,omitempty"`
			Position []float64 `json:"position,omitempty"`
			Tags     []struct {
				Group string `json:"group,omitempty"`
				ID    string `json:"id,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"tags,omitempty"`
			Title    string `json:"title,omitempty"`
			Type     string `json:"type,omitempty"`
			Vicinity string `json:"vicinity,omitempty"`
		} `json:"items,omitempty"`
		Next string `json:"next,omitempty"`
	} `json:"results,omitempty"`
	Search struct {
		Context struct {
			Href     string `json:"href,omitempty"`
			Location struct {
				Address struct {
					City        string `json:"city,omitempty"`
					Country     string `json:"country,omitempty"`
					CountryCode string `json:"countryCode,omitempty"`
					County      string `json:"county,omitempty"`
					District    string `json:"district,omitempty"`
					House       string `json:"house,omitempty"`
					PostalCode  string `json:"postalCode,omitempty"`
					StateCode   string `json:"stateCode,omitempty"`
					Street      string `json:"street,omitempty"`
					Text        string `json:"text,omitempty"`
				} `json:"address,omitempty"`
				Position []float64 `json:"position,omitempty"`
			} `json:"location,omitempty"`
			Type string `json:"type,omitempty"`
		} `json:"context,omitempty"`
	} `json:"search,omitempty"`
}

// HerePlaceDetailsResponse represent detailed info about specific place from HERE API
type HerePlaceDetailsResponse struct {
	Name             string `json:"name,omitempty"`
	PlaceID          string `json:"placeId,omitempty"`
	View             string `json:"view,omitempty"`
	AlternativeNames []struct {
		Name     string `json:"name,omitempty"`
		Language string `json:"language,omitempty"`
	} `json:"alternativeNames,omitempty"`
	Location struct {
		Position []float64 `json:"position,omitempty"`
		Address  struct {
			Text        string `json:"text,omitempty"`
			House       string `json:"house,omitempty"`
			Street      string `json:"street,omitempty"`
			PostalCode  string `json:"postalCode,omitempty"`
			District    string `json:"district,omitempty"`
			City        string `json:"city,omitempty"`
			State       string `json:"state,omitempty"`
			Country     string `json:"country,omitempty"`
			CountryCode string `json:"countryCode,omitempty"`
		} `json:"address,omitempty"`
		Access []struct {
			Position   []float64 `json:"position,omitempty"`
			AccessType string    `json:"accessType,omitempty"`
		} `json:"access,omitempty"`
	} `json:"location,omitempty"`
	Contacts struct {
		Phone []struct {
			Value string `json:"value,omitempty"`
			Label string `json:"label,omitempty"`
		} `json:"phone,omitempty"`
		Website []struct {
			Value string `json:"value,omitempty"`
			Label string `json:"label,omitempty"`
		} `json:"website,omitempty"`
	} `json:"contacts,omitempty"`
	Categories []struct {
		ID     string `json:"id,omitempty"`
		Title  string `json:"title,omitempty"`
		Href   string `json:"href,omitempty"`
		Type   string `json:"type,omitempty"`
		System string `json:"system,omitempty"`
		Icon   string `json:"icon,omitempty"`
	} `json:"categories,omitempty"`
	Tags []struct {
		ID    string `json:"id,omitempty"`
		Title string `json:"title,omitempty"`
		Group string `json:"group,omitempty"`
	} `json:"tags,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Media struct {
		Images struct {
			Available int `json:"available,omitempty"`
			//Items     []interface{} `json:"items,omitempty"`
		} `json:"images,omitempty"`
		Reviews struct {
			Available int `json:"available,omitempty"`
			//Items     []interface{} `json:"items,omitempty"`
		} `json:"reviews,omitempty"`
		Ratings struct {
			Available int `json:"available,omitempty"`
			//Items     []interface{} `json:"items,omitempty"`
		} `json:"ratings,omitempty"`
	} `json:"media,omitempty"`
	Extended struct {
		OpeningHours struct {
			Text       string `json:"text,omitempty"`
			Label      string `json:"label,omitempty"`
			IsOpen     bool   `json:"isOpen,omitempty"`
			Structured []struct {
				Start      string `json:"start,omitempty"`
				Duration   string `json:"duration,omitempty"`
				Recurrence string `json:"recurrence,omitempty"`
			} `json:"structured,omitempty"`
		} `json:"openingHours,omitempty"`
	} `json:"extended,omitempty"`
	Related struct {
		Recommended struct {
			Title string `json:"title,omitempty"`
			Href  string `json:"href,omitempty"`
			Type  string `json:"type,omitempty"`
		} `json:"recommended,omitempty"`
		PublicTransport struct {
			Title string `json:"title,omitempty"`
			Href  string `json:"href,omitempty"`
			Type  string `json:"type,omitempty"`
		} `json:"public-transport,omitempty"`
	} `json:"related,omitempty"`
	Report struct {
		Title string `json:"title,omitempty"`
		Href  string `json:"href,omitempty"`
		Type  string `json:"type,omitempty"`
	} `json:"report,omitempty"`
}

// HereShortURLResponse represent short version of URL
type HereShortURLResponse struct {
	URL     string `json:"url,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// HereGeocodingResponse represent response from HERE Geocoding API
type HereGeocodingResponse struct {
	Response struct {
		MetaInfo struct {
			Timestamp string `json:"Timestamp,omitempty"`
		} `json:"MetaInfo,omitempty"`
		View []struct {
			Type   string `json:"_type,omitempty"`
			ViewID int    `json:"ViewId,omitempty"`
			Result []struct {
				Relevance    float64 `json:"Relevance,omitempty"`
				MatchLevel   string  `json:"MatchLevel,omitempty"`
				MatchQuality struct {
					City float64 `json:"City,omitempty"`
				} `json:"MatchQuality,omitempty"`
				Location struct {
					LocationID      string `json:"LocationId,omitempty"`
					LocationType    string `json:"LocationType,omitempty"`
					DisplayPosition struct {
						Latitude  float64 `json:"Latitude,omitempty"`
						Longitude float64 `json:"Longitude,omitempty"`
					} `json:"DisplayPosition,omitempty"`
					NavigationPosition []struct {
						Latitude  float64 `json:"Latitude,omitempty"`
						Longitude float64 `json:"Longitude,omitempty"`
					} `json:"NavigationPosition,omitempty"`
					MapView struct {
						TopLeft struct {
							Latitude  float64 `json:"Latitude,omitempty"`
							Longitude float64 `json:"Longitude,omitempty"`
						} `json:"TopLeft,omitempty"`
						BottomRight struct {
							Latitude  float64 `json:"Latitude,omitempty"`
							Longitude float64 `json:"Longitude,omitempty"`
						} `json:"BottomRight,omitempty"`
					} `json:"MapView,omitempty"`
					Address struct {
						Label          string `json:"Label,omitempty"`
						Country        string `json:"Country,omitempty"`
						State          string `json:"State,omitempty"`
						County         string `json:"County,omitempty"`
						City           string `json:"City,omitempty"`
						AdditionalData []struct {
							Value string `json:"value,omitempty"`
							Key   string `json:"key,omitempty"`
						} `json:"AdditionalData,omitempty"`
					} `json:"Address,omitempty"`
				} `json:"Location,omitempty"`
			} `json:"Result,omitempty"`
		} `json:"View,omitempty"`
	} `json:"Response,omitempty"`
}
