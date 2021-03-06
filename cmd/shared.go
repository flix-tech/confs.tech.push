package cmd

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/urfave/cli.v1"

	"github.com/flix-tech/confs.tech.push/confs"
)

func validateTopicArgument(topic string) (string, error) {
	if topic == "" {
		return "", errors.New("Please provide conference topic")
	}

	match, _ := regexp.MatchString("^[a-z\\-]+$", topic)
	if !match {
		return "", errors.New("Invalid conference topic")
	}

	return topic, nil
}

func wrapAction(action func(topic string, conferences []confs.Conference, c *cli.Context) error) func(c *cli.Context) error {
	return func (c *cli.Context) error {
		topic, err := validateTopicArgument(c.Args().Get(0))
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		// Fetch conference data
		conferences, err := confs.GetConferences(topic)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		conferences = confs.FilterConferences(conferences,
			confs.NewIsInFutureTest(),
			confs.NewCFPFinishedTest(c.GlobalBool("cfp-finished")),
			confs.NewIsNotInBlacklistedCountryTest(c.GlobalStringSlice("countries-blacklist")),
		)

		err = action(topic, conferences, c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	}
}

func formatDateRange(c confs.Conference) string {
	dateRange := c.StartDate
	if c.StartDate != c.EndDate {
		dateRange = fmt.Sprintf("%s — %s", c.StartDate, c.EndDate)
	}
	return dateRange
}

func formatLocation(c confs.Conference) string {
	flags := map[string]string{
		"Ascension Island": "🇦🇨",
		"Andorra": "🇦🇩",
		"United Arab Emirates": "🇦🇪",
		"Afghanistan": "🇦🇫",
		"Antigua & Barbuda": "🇦🇬",
		"Anguilla": "🇦🇮",
		"Albania": "🇦🇱",
		"Armenia": "🇦🇲",
		"Angola": "🇦🇴",
		"Antarctica": "🇦🇶",
		"Argentina": "🇦🇷",
		"American Samoa": "🇦🇸",
		"Austria": "🇦🇹",
		"Australia": "🇦🇺",
		"Aruba": "🇦🇼",
		"Åland Islands": "🇦🇽",
		"Azerbaijan": "🇦🇿",
		"Bosnia & Herzegovina": "🇧🇦",
		"Barbados": "🇧🇧",
		"Bangladesh": "🇧🇩",
		"Belgium": "🇧🇪",
		"Burkina Faso": "🇧🇫",
		"Bulgaria": "🇧🇬",
		"Bahrain": "🇧🇭",
		"Burundi": "🇧🇮",
		"Benin": "🇧🇯",
		"St. Barthélemy": "🇧🇱",
		"Bermuda": "🇧🇲",
		"Brunei": "🇧🇳",
		"Bolivia": "🇧🇴",
		"Caribbean Netherlands": "🇧🇶",
		"Brazil": "🇧🇷",
		"Bahamas": "🇧🇸",
		"Bhutan": "🇧🇹",
		"Bouvet Island": "🇧🇻",
		"Botswana": "🇧🇼",
		"Belarus": "🇧🇾",
		"Belize": "🇧🇿",
		"Canada": "🇨🇦",
		"Cocos (Keeling) Islands": "🇨🇨",
		"Congo - Kinshasa": "🇨🇩",
		"Central African Republic": "🇨🇫",
		"Congo - Brazzaville": "🇨🇬",
		"Switzerland": "🇨🇭",
		"Côte d’Ivoire": "🇨🇮",
		"Cook Islands": "🇨🇰",
		"Chile": "🇨🇱",
		"Cameroon": "🇨🇲",
		"China": "🇨🇳",
		"Colombia": "🇨🇴",
		"Clipperton Island": "🇨🇵",
		"Costa Rica": "🇨🇷",
		"Cuba": "🇨🇺",
		"Cape Verde": "🇨🇻",
		"Curaçao": "🇨🇼",
		"Christmas Island": "🇨🇽",
		"Cyprus": "🇨🇾",
		"Czechia": "🇨🇿",
		"Czech Republic": "🇨🇿",
		"Germany": "🇩🇪",
		"Deutschland": "🇩🇪",
		"Diego Garcia": "🇩🇬",
		"Djibouti": "🇩🇯",
		"Denmark": "🇩🇰",
		"Dominica": "🇩🇲",
		"Dominican Republic": "🇩🇴",
		"Algeria": "🇩🇿",
		"Ceuta & Melilla": "🇪🇦",
		"Ecuador": "🇪🇨",
		"Estonia": "🇪🇪",
		"Egypt": "🇪🇬",
		"Western Sahara": "🇪🇭",
		"Eritrea": "🇪🇷",
		"Spain": "🇪🇸",
		"Ethiopia": "🇪🇹",
		"European Union": "🇪🇺",
		"Finland": "🇫🇮",
		"Fiji": "🇫🇯",
		"Falkland Islands": "🇫🇰",
		"Micronesia": "🇫🇲",
		"Faroe Islands": "🇫🇴",
		"France": "🇫🇷",
		"Gabon": "🇬🇦",
		"United Kingdom": "🇬🇧",
		"U.K.": "🇬🇧",
		"Grenada": "🇬🇩",
		"Georgia": "🇬🇪",
		"French Guiana": "🇬🇫",
		"Guernsey": "🇬🇬",
		"Ghana": "🇬🇭",
		"Gibraltar": "🇬🇮",
		"Greenland": "🇬🇱",
		"Gambia": "🇬🇲",
		"Guinea": "🇬🇳",
		"Guadeloupe": "🇬🇵",
		"Equatorial Guinea": "🇬🇶",
		"Greece": "🇬🇷",
		"South Georgia & South Sandwich Islands": "🇬🇸",
		"Guatemala": "🇬🇹",
		"Guam": "🇬🇺",
		"Guinea-Bissau": "🇬🇼",
		"Guyana": "🇬🇾",
		"Hong Kong SAR China": "🇭🇰",
		"Heard & McDonald Islands": "🇭🇲",
		"Honduras": "🇭🇳",
		"Croatia": "🇭🇷",
		"Haiti": "🇭🇹",
		"Hungary": "🇭🇺",
		"Canary Islands": "🇮🇨",
		"Indonesia": "🇮🇩",
		"Ireland": "🇮🇪",
		"Israel": "🇮🇱",
		"Isle of Man": "🇮🇲",
		"India": "🇮🇳",
		"British Indian Ocean Territory": "🇮🇴",
		"Iraq": "🇮🇶",
		"Iran": "🇮🇷",
		"Iceland": "🇮🇸",
		"Italy": "🇮🇹",
		"Jersey": "🇯🇪",
		"Jamaica": "🇯🇲",
		"Jordan": "🇯🇴",
		"Japan": "🇯🇵",
		"Kenya": "🇰🇪",
		"Kyrgyzstan": "🇰🇬",
		"Cambodia": "🇰🇭",
		"Kiribati": "🇰🇮",
		"Comoros": "🇰🇲",
		"St. Kitts & Nevis": "🇰🇳",
		"North Korea": "🇰🇵",
		"South Korea": "🇰🇷",
		"Kuwait": "🇰🇼",
		"Cayman Islands": "🇰🇾",
		"Kazakhstan": "🇰🇿",
		"Laos": "🇱🇦",
		"Lebanon": "🇱🇧",
		"St. Lucia": "🇱🇨",
		"Liechtenstein": "🇱🇮",
		"Sri Lanka": "🇱🇰",
		"Liberia": "🇱🇷",
		"Lesotho": "🇱🇸",
		"Lithuania": "🇱🇹",
		"Luxembourg": "🇱🇺",
		"Latvia": "🇱🇻",
		"Libya": "🇱🇾",
		"Morocco": "🇲🇦",
		"Monaco": "🇲🇨",
		"Moldova": "🇲🇩",
		"Montenegro": "🇲🇪",
		"St. Martin": "🇲🇫",
		"Madagascar": "🇲🇬",
		"Marshall Islands": "🇲🇭",
		"North Macedonia": "🇲🇰",
		"Mali": "🇲🇱",
		"Myanmar (Burma)": "🇲🇲",
		"Mongolia": "🇲🇳",
		"Macau Sar China": "🇲🇴",
		"Northern Mariana Islands": "🇲🇵",
		"Martinique": "🇲🇶",
		"Mauritania": "🇲🇷",
		"Montserrat": "🇲🇸",
		"Malta": "🇲🇹",
		"Mauritius": "🇲🇺",
		"Maldives": "🇲🇻",
		"Malawi": "🇲🇼",
		"Mexico": "🇲🇽",
		"Malaysia": "🇲🇾",
		"Mozambique": "🇲🇿",
		"Namibia": "🇳🇦",
		"New Caledonia": "🇳🇨",
		"Niger": "🇳🇪",
		"Norfolk Island": "🇳🇫",
		"Nigeria": "🇳🇬",
		"Nicaragua": "🇳🇮",
		"Netherlands": "🇳🇱",
		"Norway": "🇳🇴",
		"Nepal": "🇳🇵",
		"Nauru": "🇳🇷",
		"Niue": "🇳🇺",
		"New Zealand": "🇳🇿",
		"Oman": "🇴🇲",
		"Panama": "🇵🇦",
		"Peru": "🇵🇪",
		"French Polynesia": "🇵🇫",
		"Papua New Guinea": "🇵🇬",
		"Philippines": "🇵🇭",
		"Pakistan": "🇵🇰",
		"Poland": "🇵🇱",
		"St. Pierre & Miquelon": "🇵🇲",
		"Pitcairn Islands": "🇵🇳",
		"Puerto Rico": "🇵🇷",
		"Palestinian Territories": "🇵🇸",
		"Portugal": "🇵🇹",
		"Palau": "🇵🇼",
		"Paraguay": "🇵🇾",
		"Qatar": "🇶🇦",
		"Réunion": "🇷🇪",
		"Romania": "🇷🇴",
		"Serbia": "🇷🇸",
		"Russia": "🇷🇺",
		"Rwanda": "🇷🇼",
		"Saudi Arabia": "🇸🇦",
		"Solomon Islands": "🇸🇧",
		"Seychelles": "🇸🇨",
		"Sudan": "🇸🇩",
		"Sweden": "🇸🇪",
		"Singapore": "🇸🇬",
		"St. Helena": "🇸🇭",
		"Slovenia": "🇸🇮",
		"Svalbard & Jan Mayen": "🇸🇯",
		"Slovakia": "🇸🇰",
		"Sierra Leone": "🇸🇱",
		"San Marino": "🇸🇲",
		"Senegal": "🇸🇳",
		"Somalia": "🇸🇴",
		"Suriname": "🇸🇷",
		"South Sudan": "🇸🇸",
		"São Tomé & Príncipe": "🇸🇹",
		"El Salvador": "🇸🇻",
		"Sint Maarten": "🇸🇽",
		"Syria": "🇸🇾",
		"Swaziland": "🇸🇿",
		"Tristan Da Cunha": "🇹🇦",
		"Turks & Caicos Islands": "🇹🇨",
		"Chad": "🇹🇩",
		"French Southern Territories": "🇹🇫",
		"Togo": "🇹🇬",
		"Thailand": "🇹🇭",
		"Tajikistan": "🇹🇯",
		"Tokelau": "🇹🇰",
		"Timor-Leste": "🇹🇱",
		"Turkmenistan": "🇹🇲",
		"Tunisia": "🇹🇳",
		"Tonga": "🇹🇴",
		"Turkey": "🇹🇷",
		"Trinidad & Tobago": "🇹🇹",
		"Tuvalu": "🇹🇻",
		"Taiwan": "🇹🇼",
		"Tanzania": "🇹🇿",
		"Ukraine": "🇺🇦",
		"Uganda": "🇺🇬",
		"U.S. Outlying Islands": "🇺🇲",
		"United States": "🇺🇸",
		"U.S.A.": "🇺🇸",
		"USA": "🇺🇸",
		"Uruguay": "🇺🇾",
		"Uzbekistan": "🇺🇿",
		"Vatican City": "🇻🇦",
		"St. Vincent & Grenadines": "🇻🇨",
		"Venezuela": "🇻🇪",
		"British Virgin Islands": "🇻🇬",
		"U.S. Virgin Islands": "🇻🇮",
		"Vietnam": "🇻🇳",
		"Vanuatu": "🇻🇺",
		"Wallis & Futuna": "🇼🇫",
		"Samoa": "🇼🇸",
		"Kosovo": "🇽🇰",
		"Yemen": "🇾🇪",
		"Mayotte": "🇾🇹",
		"South Africa": "🇿🇦",
		"Zambia": "🇿🇲",
		"Zimbabwe": "🇿🇼",
		"England": "🏴󠁧󠁢󠁥󠁮󠁧󠁿",
		"Scotland": "🏴󠁧󠁢󠁳󠁣󠁴󠁿",
	}

	location := fmt.Sprintf("%s, %s", c.City, c.Country)

	flag, flagFound := flags[c.Country]
	if (flagFound) {
		location += " " + flag
	}

	return location
}
