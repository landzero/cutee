package types

// TAccount twitter account object
type TAccount struct {
	ID         string `json:"id_str"`
	ScreenName string `json:"screen_name"`
	AvatarURL  string `json:"profile_image_url_https"`
}
