package geoutil

import (
	"strings"

	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

const (
	Afghanistan                            = `af`
	Albania                                = `al`
	Algeria                                = `dz`
	AmericanSamoa                          = `as`
	Andorra                                = `ad`
	Angola                                 = `ao`
	Anguilla                               = `ai`
	Antarctica                             = `aq`
	AntiguaAndBarbuda                      = `ag`
	Argentina                              = `ar`
	Armenia                                = `am`
	Aruba                                  = `aw`
	Australia                              = `au`
	Austria                                = `at`
	Azerbaijan                             = `az`
	Bahamas                                = `bs`
	Bahrain                                = `bh`
	Bangladesh                             = `bd`
	Barbados                               = `bb`
	Belarus                                = `by`
	Belgium                                = `be`
	Belize                                 = `bz`
	Benin                                  = `bj`
	Bermuda                                = `bm`
	Bhutan                                 = `bt`
	Bolivia                                = `bo`
	BosniaAndHerzegovina                   = `ba`
	Botswana                               = `bw`
	BouvetIsland                           = `bv`
	Brazil                                 = `br`
	BritishIndianOceanTerritory            = `io`
	BritishVirginIslands                   = `vg`
	Brunei                                 = `bn`
	Bulgaria                               = `bg`
	BurkinaFaso                            = `bf`
	Burma                                  = `mm`
	Burundi                                = `bi`
	Cambodia                               = `kh`
	Cameroon                               = `cm`
	Canada                                 = `ca`
	CapeVerde                              = `cv`
	CaymanIslands                          = `ky`
	CentralAfricanRepublic                 = `cf`
	Chad                                   = `td`
	Chile                                  = `cl`
	China                                  = `cn`
	ChristmasIsland                        = `cx`
	CocosIslands                           = `cc`
	Colombia                               = `co`
	Comoros                                = `km`
	CongoDRC                               = `cd`
	CongoRepublic                          = `cg`
	CookIslands                            = `ck`
	CostaRica                              = `cr`
	CoteDIvoire                            = `ci`
	CôtedIvoire                            = `ci`
	Croatia                                = `hr`
	Cuba                                   = `cu`
	Cyprus                                 = `cy`
	CzechRepublic                          = `cz`
	Denmark                                = `dk`
	Djibouti                               = `dj`
	Dominica                               = `dm`
	DominicanRepublic                      = `do`
	Ecuador                                = `ec`
	Egypt                                  = `eg`
	ElSalvador                             = `sv`
	EquatorialGuinea                       = `gq`
	Eritrea                                = `er`
	Estonia                                = `ee`
	Ethiopia                               = `et`
	FalklandIslands                        = `fk`
	FaroeIslands                           = `fo`
	Fiji                                   = `fj`
	Finland                                = `fi`
	France                                 = `fr`
	FrenchGuiana                           = `gf`
	FrenchPolynesia                        = `pf`
	FrenchSouthernTerritories              = `tf`
	Gabon                                  = `ga`
	Gambia                                 = `gm`
	GazaStrip                              = `gz`
	Georgia                                = `ge`
	Germany                                = `de`
	Ghana                                  = `gh`
	Gibraltar                              = `gi`
	GreatBritain                           = `gb`
	Greece                                 = `gr`
	Greenland                              = `gl`
	Grenada                                = `gd`
	Guadeloupe                             = `gp`
	Guam                                   = `gu`
	Guatemala                              = `gt`
	Guernsey                               = `gg`
	Guinea                                 = `gn`
	GuineaBissau                           = `gw`
	Guyana                                 = `gy`
	Haiti                                  = `ht`
	HeardIslandAndMcDonaldIslands          = `hm`
	Honduras                               = `hn`
	HongKong                               = `hk`
	Hungary                                = `hu`
	Iceland                                = `is`
	India                                  = `in`
	Indonesia                              = `id`
	Iran                                   = `ir`
	Iraq                                   = `iq`
	Ireland                                = `ie`
	IslasMalvinas                          = `fk`
	IsleOfMan                              = `im`
	Israel                                 = `il`
	Italy                                  = `it`
	Jamaica                                = `jm`
	JanMayen                               = `sj`
	Japan                                  = `jp`
	Jersey                                 = `je`
	Jordan                                 = `jo`
	Kazakhstan                             = `kz`
	Kenya                                  = `ke`
	Kiribati                               = `ki`
	Kosovo                                 = `xk`
	Kuwait                                 = `kw`
	Kyrgyzstan                             = `kg`
	Laos                                   = `la`
	Latvia                                 = `lv`
	Lebanon                                = `lb`
	Lesotho                                = `ls`
	Liberia                                = `lr`
	Libya                                  = `ly`
	Liechtenstein                          = `li`
	Lithuania                              = `lt`
	Luxembourg                             = `lu`
	Macau                                  = `mo`
	Macedonia                              = `mk`
	Madagascar                             = `mg`
	Malawi                                 = `mw`
	Malaysia                               = `my`
	Maldives                               = `mv`
	Mali                                   = `ml`
	Malta                                  = `mt`
	MarshallIslands                        = `mh`
	Martinique                             = `mq`
	Mauritania                             = `mr`
	Mauritius                              = `mu`
	Mayotte                                = `yt`
	Mexico                                 = `mx`
	Micronesia                             = `fm`
	Miquelon                               = `pm`
	Moldova                                = `md`
	Monaco                                 = `mc`
	Mongolia                               = `mn`
	Montenegro                             = `me`
	Montserrat                             = `ms`
	Morocco                                = `ma`
	Mozambique                             = `mz`
	Myanmar                                = `mm`
	Namibia                                = `na`
	Nauru                                  = `nr`
	Nepal                                  = `np`
	Netherlands                            = `nl`
	NetherlandsAntilles                    = `an`
	NewCaledonia                           = `nc`
	NewZealand                             = `nz`
	Nicaragua                              = `ni`
	Niger                                  = `ne`
	Nigeria                                = `ng`
	Niue                                   = `nu`
	NorfolkIsland                          = `nf`
	NorthernMarianaIslands                 = `mp`
	NorthKorea                             = `kp`
	Norway                                 = `no`
	Oman                                   = `om`
	Pakistan                               = `pk`
	Palau                                  = `pw`
	Palestine                              = `ps`
	Panama                                 = `pa`
	PapuaNewGuinea                         = `pg`
	Paraguay                               = `py`
	Peru                                   = `pe`
	Philippines                            = `ph`
	PitcairnIslands                        = `pn`
	Poland                                 = `pl`
	Portugal                               = `pt`
	PuertoRico                             = `pr`
	Qatar                                  = `qa`
	Reunion                                = `re`
	Réunion                                = `re`
	Romania                                = `ro`
	Russia                                 = `ru`
	Rwanda                                 = `rw`
	SaintHelena                            = `sh`
	SaintKittsAndNevis                     = `kn`
	SaintLucia                             = `lc`
	SaintPierre                            = `pm`
	SaintPierreAndMiquelon                 = `pm`
	SaintVincent                           = `vc`
	SaintVincentAndTheGrenadines           = `vc`
	Samoa                                  = `ws`
	SanMarino                              = `sm`
	SaoTomeAndPrincipe                     = `st`
	SãoToméAndPríncipe                     = `st`
	SaudiArabia                            = `sa`
	Senegal                                = `sn`
	Serbia                                 = `rs`
	Seychelles                             = `sc`
	SierraLeone                            = `sl`
	Singapore                              = `sg`
	Slovakia                               = `sk`
	Slovenia                               = `si`
	SolomonIslands                         = `sb`
	Somalia                                = `so`
	SouthAfrica                            = `za`
	SouthGeorgia                           = `gs`
	SouthGeorgiaAndTheSouthSandwichIslands = `gs`
	SouthKorea                             = `kr`
	SouthSandwichIslands                   = `gs`
	SouthSudan                             = `ss`
	SovietUnion                            = `su`
	Spain                                  = `es`
	SriLanka                               = `lk`
	Sudan                                  = `sd`
	Suriname                               = `sr`
	Svalbard                               = `sj`
	SvalbardAndJanMayen                    = `sj`
	Swaziland                              = `sz`
	Sweden                                 = `se`
	Switzerland                            = `ch`
	Syria                                  = `sy`
	Taiwan                                 = `tw`
	Tajikistan                             = `tj`
	Tanzania                               = `tz`
	Thailand                               = `th`
	TheGrenadines                          = `vc`
	TimorLeste                             = `tl`
	Togo                                   = `tg`
	Tokelau                                = `tk`
	Tonga                                  = `to`
	TrinidadAndTobago                      = `tt`
	Tunisia                                = `tn`
	Turkey                                 = `tr`
	Turkmenistan                           = `tm`
	TurksAndCaicosIslands                  = `tc`
	Tuvalu                                 = `tv`
	Uganda                                 = `ug`
	UK                                     = `gb`
	Ukraine                                = `ua`
	UnitedArabEmirates                     = `ae`
	UnitedKingdom                          = `gb`
	UnitedStates                           = `us`
	UnitedStatesMinorOutlyingIslands       = `um`
	UnitedStatesVirginIslands              = `vi`
	Uruguay                                = `uy`
	Uzbekistan                             = `uz`
	Vanuatu                                = `vu`
	VaticanCity                            = `va`
	Venezuela                              = `ve`
	Vietnam                                = `vn`
	WallisAndFutuna                        = `wf`
	WesternSahara                          = `eh`
	Yemen                                  = `ye`
	Zambia                                 = `zm`
	Zimbabwe                               = `zw`
)

type Country struct {
	Name      string
	Code      string
	CCTLD     string
	Latitude  float64
	Longitude float64
}

func (self Country) IsValid() bool {
	return (self.Name != ``) && (self.Code != ``)
}

type CountryData map[string]Country

func (self CountryData) Get(nameOrCode any) Country {
	var ia2 = strings.ToLower(typeutil.String(nameOrCode))

	if c, ok := self[ia2]; ok && c.IsValid() {
		return c
	}

	for k, v := range self {
		if stringutil.SoftEqual(v.Code, nameOrCode) {
			return v
		} else if stringutil.SoftEqual(k, nameOrCode) {
			return v
		}
	}

	return Country{}
}

// A data structure that exposes standard information about countries of the world, keyed on their ISO3166-1 Alpha2 abbreviation (e.g.: "us", "gb", "de", etc.)
var Countries = CountryData{
	Andorra: {
		Name:      `Andorra`,
		Code:      Andorra,
		Latitude:  42.546245,
		Longitude: 1.601554,
		CCTLD:     `.ad`,
	},
	UnitedArabEmirates: {
		Name:      `United Arab Emirates`,
		Code:      UnitedArabEmirates,
		Latitude:  23.424076,
		Longitude: 53.847818,
		CCTLD:     `.ae`,
	},
	Afghanistan: {
		Name:      `Afghanistan`,
		Code:      Afghanistan,
		Latitude:  33.93911,
		Longitude: 67.709953,
		CCTLD:     `.af`,
	},
	AntiguaAndBarbuda: {
		Name:      `Antigua and Barbuda`,
		Code:      AntiguaAndBarbuda,
		Latitude:  17.060816,
		Longitude: -61.796428,
		CCTLD:     `.ag`,
	},
	Anguilla: {
		Name:      `Anguilla`,
		Code:      Anguilla,
		Latitude:  18.220554,
		Longitude: -63.068615,
		CCTLD:     `.ai`,
	},
	Albania: {
		Name:      `Albania`,
		Code:      Albania,
		Latitude:  41.153332,
		Longitude: 20.168331,
		CCTLD:     `.al`,
	},
	Armenia: {
		Name:      `Armenia`,
		Code:      Armenia,
		Latitude:  40.069099,
		Longitude: 45.038189,
		CCTLD:     `.am`,
	},
	NetherlandsAntilles: {
		Name:      `Netherlands Antilles`,
		Code:      NetherlandsAntilles,
		Latitude:  12.226079,
		Longitude: -69.060087,
	},
	Angola: {
		Name:      `Angola`,
		Code:      Angola,
		Latitude:  -11.202692,
		Longitude: 17.873887,
		CCTLD:     `.ao`,
	},
	Antarctica: {
		Name:      `Antarctica`,
		Code:      Antarctica,
		Latitude:  -75.250973,
		Longitude: -0.071389,
		CCTLD:     `.aq`,
	},
	Argentina: {
		Name:      `Argentina`,
		Code:      Argentina,
		Latitude:  -38.416097,
		Longitude: -63.616672,
		CCTLD:     `.ar`,
	},
	AmericanSamoa: {
		Name:      `American Samoa`,
		Code:      AmericanSamoa,
		Latitude:  -14.270972,
		Longitude: -170.132217,
		CCTLD:     `.as`,
	},
	Austria: {
		Name:      `Austria`,
		Code:      Austria,
		Latitude:  47.516231,
		Longitude: 14.550072,
		CCTLD:     `.at`,
	},
	Australia: {
		Name:      `Australia`,
		Code:      Australia,
		Latitude:  -25.274398,
		Longitude: 133.775136,
		CCTLD:     `.au`,
	},
	Aruba: {
		Name:      `Aruba`,
		Code:      Aruba,
		Latitude:  12.52111,
		Longitude: -69.968338,
		CCTLD:     `.aw`,
	},
	Azerbaijan: {
		Name:      `Azerbaijan`,
		Code:      Azerbaijan,
		Latitude:  40.143105,
		Longitude: 47.576927,
		CCTLD:     `.az`,
	},
	BosniaAndHerzegovina: {
		Name:      `Bosnia and Herzegovina`,
		Code:      BosniaAndHerzegovina,
		Latitude:  43.915886,
		Longitude: 17.679076,
		CCTLD:     `.ba`,
	},
	Barbados: {
		Name:      `Barbados`,
		Code:      Barbados,
		Latitude:  13.193887,
		Longitude: -59.543198,
		CCTLD:     `.bb`,
	},
	Bangladesh: {
		Name:      `Bangladesh`,
		Code:      Bangladesh,
		Latitude:  23.684994,
		Longitude: 90.356331,
		CCTLD:     `.bd`,
	},
	Belgium: {
		Name:      `Belgium`,
		Code:      Belgium,
		Latitude:  50.503887,
		Longitude: 4.469936,
		CCTLD:     `.be`,
	},
	BurkinaFaso: {
		Name:      `Burkina Faso`,
		Code:      BurkinaFaso,
		Latitude:  12.238333,
		Longitude: -1.561593,
		CCTLD:     `.bf`,
	},
	Bulgaria: {
		Name:      `Bulgaria`,
		Code:      Bulgaria,
		Latitude:  42.733883,
		Longitude: 25.48583,
		CCTLD:     `.bg`,
	},
	Bahrain: {
		Name:      `Bahrain`,
		Code:      Bahrain,
		Latitude:  25.930414,
		Longitude: 50.637772,
		CCTLD:     `.bh`,
	},
	Burundi: {
		Name:      `Burundi`,
		Code:      Burundi,
		Latitude:  -3.373056,
		Longitude: 29.918886,
		CCTLD:     `.bi`,
	},
	Benin: {
		Name:      `Benin`,
		Code:      Benin,
		Latitude:  9.30769,
		Longitude: 2.315834,
		CCTLD:     `.bj`,
	},
	Bermuda: {
		Name:      `Bermuda`,
		Code:      Bermuda,
		Latitude:  32.321384,
		Longitude: -64.75737,
		CCTLD:     `.bm`,
	},
	Brunei: {
		Name:      `Brunei`,
		Code:      Brunei,
		Latitude:  4.535277,
		Longitude: 114.727669,
		CCTLD:     `.bn`,
	},
	Bolivia: {
		Name:      `Bolivia`,
		Code:      Bolivia,
		Latitude:  -16.290154,
		Longitude: -63.588653,
		CCTLD:     `.bo`,
	},
	Brazil: {
		Name:      `Brazil`,
		Code:      Brazil,
		Latitude:  -14.235004,
		Longitude: -51.92528,
		CCTLD:     `.br`,
	},
	Bahamas: {
		Name:      `Bahamas`,
		Code:      Bahamas,
		Latitude:  25.03428,
		Longitude: -77.39628,
		CCTLD:     `.bs`,
	},
	Bhutan: {
		Name:      `Bhutan`,
		Code:      Bhutan,
		Latitude:  27.514162,
		Longitude: 90.433601,
		CCTLD:     `.bt`,
	},
	BouvetIsland: {
		Name:      `Bouvet Island`,
		Code:      BouvetIsland,
		Latitude:  -54.423199,
		Longitude: 3.413194,
	},
	Botswana: {
		Name:      `Botswana`,
		Code:      Botswana,
		Latitude:  -22.328474,
		Longitude: 24.684866,
		CCTLD:     `.bw`,
	},
	Belarus: {
		Name:      `Belarus`,
		Code:      Belarus,
		Latitude:  53.709807,
		Longitude: 27.953389,
		CCTLD:     `.by`,
	},
	Belize: {
		Name:      `Belize`,
		Code:      Belize,
		Latitude:  17.189877,
		Longitude: -88.49765,
		CCTLD:     `.bz`,
	},
	Canada: {
		Name:      `Canada`,
		Code:      Canada,
		Latitude:  56.130366,
		Longitude: -106.346771,
		CCTLD:     `.ca`,
	},
	CocosIslands: {
		Name:      `Cocos Islands`,
		Code:      CocosIslands,
		Latitude:  -12.164165,
		Longitude: 96.870956,
		CCTLD:     `.cc`,
	},
	CongoDRC: {
		Name:      `Democratic Republic of Congo`,
		Code:      CongoDRC,
		Latitude:  -4.038333,
		Longitude: 21.758664,
		CCTLD:     `.cd`,
	},
	CentralAfricanRepublic: {
		Name:      `Central African Republic`,
		Code:      CentralAfricanRepublic,
		Latitude:  6.611111,
		Longitude: 20.939444,
		CCTLD:     `.cf`,
	},
	CongoRepublic: {
		Name:      `Congo Republic`,
		Code:      CongoRepublic,
		Latitude:  -0.228021,
		Longitude: 15.827659,
		CCTLD:     `.cg`,
	},
	Switzerland: {
		Name:      `Switzerland`,
		Code:      Switzerland,
		Latitude:  46.818188,
		Longitude: 8.227512,
		CCTLD:     `.ch`,
	},
	CoteDIvoire: {
		Name:      `Côte d'Ivoire`,
		Code:      CoteDIvoire,
		Latitude:  7.539989,
		Longitude: -5.54708,
		CCTLD:     `.ci`,
	},
	CookIslands: {
		Name:      `Cook Islands`,
		Code:      CookIslands,
		Latitude:  -21.236736,
		Longitude: -159.777671,
		CCTLD:     `.ck`,
	},
	Chile: {
		Name:      `Chile`,
		Code:      Chile,
		Latitude:  -35.675147,
		Longitude: -71.542969,
		CCTLD:     `.cl`,
	},
	Cameroon: {
		Name:      `Cameroon`,
		Code:      Cameroon,
		Latitude:  7.369722,
		Longitude: 12.354722,
		CCTLD:     `.cm`,
	},
	China: {
		Name:      `China`,
		Code:      China,
		Latitude:  35.86166,
		Longitude: 104.195397,
		CCTLD:     `.cn`,
	},
	Colombia: {
		Name:      `Colombia`,
		Code:      Colombia,
		Latitude:  4.570868,
		Longitude: -74.297333,
		CCTLD:     `.co`,
	},
	CostaRica: {
		Name:      `Costa Rica`,
		Code:      CostaRica,
		Latitude:  9.748917,
		Longitude: -83.753428,
		CCTLD:     `.cr`,
	},
	Cuba: {
		Name:      `Cuba`,
		Code:      Cuba,
		Latitude:  21.521757,
		Longitude: -77.781167,
		CCTLD:     `.cu`,
	},
	CapeVerde: {
		Name:      `Cape Verde`,
		Code:      CapeVerde,
		Latitude:  16.002082,
		Longitude: -24.013197,
		CCTLD:     `.cv`,
	},
	ChristmasIsland: {
		Name:      `Christmas Island`,
		Code:      ChristmasIsland,
		Latitude:  -10.447525,
		Longitude: 105.690449,
		CCTLD:     `.cx`,
	},
	Cyprus: {
		Name:      `Cyprus`,
		Code:      Cyprus,
		Latitude:  35.126413,
		Longitude: 33.429859,
		CCTLD:     `.cy`,
	},
	CzechRepublic: {
		Name:      `Czech Republic`,
		Code:      CzechRepublic,
		Latitude:  49.817492,
		Longitude: 15.472962,
		CCTLD:     `.cz`,
	},
	Germany: {
		Name:      `Germany`,
		Code:      Germany,
		Latitude:  51.165691,
		Longitude: 10.451526,
		CCTLD:     `.de`,
	},
	Djibouti: {
		Name:      `Djibouti`,
		Code:      Djibouti,
		Latitude:  11.825138,
		Longitude: 42.590275,
		CCTLD:     `.dj`,
	},
	Denmark: {
		Name:      `Denmark`,
		Code:      Denmark,
		Latitude:  56.26392,
		Longitude: 9.501785,
		CCTLD:     `.dk`,
	},
	Dominica: {
		Name:      `Dominica`,
		Code:      Dominica,
		Latitude:  15.414999,
		Longitude: -61.370976,
		CCTLD:     `.dm`,
	},
	DominicanRepublic: {
		Name:      `Dominican Republic`,
		Code:      DominicanRepublic,
		Latitude:  18.735693,
		Longitude: -70.162651,
		CCTLD:     `.do`,
	},
	Algeria: {
		Name:      `Algeria`,
		Code:      Algeria,
		Latitude:  28.033886,
		Longitude: 1.659626,
		CCTLD:     `.dz`,
	},
	Ecuador: {
		Name:      `Ecuador`,
		Code:      Ecuador,
		Latitude:  -1.831239,
		Longitude: -78.183406,
		CCTLD:     `.ec`,
	},
	Estonia: {
		Name:      `Estonia`,
		Code:      Estonia,
		Latitude:  58.595272,
		Longitude: 25.013607,
		CCTLD:     `.ee`,
	},
	Egypt: {
		Name:      `Egypt`,
		Code:      Egypt,
		Latitude:  26.820553,
		Longitude: 30.802498,
		CCTLD:     `.eg`,
	},
	WesternSahara: {
		Name:      `Western Sahara`,
		Code:      WesternSahara,
		Latitude:  24.215527,
		Longitude: -12.885834,
		CCTLD:     `.eh`,
	},
	Eritrea: {
		Name:      `Eritrea`,
		Code:      Eritrea,
		Latitude:  15.179384,
		Longitude: 39.782334,
		CCTLD:     `.er`,
	},
	Spain: {
		Name:      `Spain`,
		Code:      Spain,
		Latitude:  40.463667,
		Longitude: -3.74922,
		CCTLD:     `.es`,
	},
	Ethiopia: {
		Name:      `Ethiopia`,
		Code:      Ethiopia,
		Latitude:  9.145,
		CCTLD:     `.et`,
		Longitude: 40.489673,
	},
	Finland: {
		Name:      `Finland`,
		Code:      Finland,
		Latitude:  61.92411,
		Longitude: 25.748151,
		CCTLD:     `.fi`,
	},
	Fiji: {
		Name:      `Fiji`,
		Code:      Fiji,
		Latitude:  -16.578193,
		Longitude: 179.414413,
		CCTLD:     `.fj`,
	},
	FalklandIslands: {
		Name:      `Falkland Islands`,
		Code:      FalklandIslands,
		Latitude:  -51.796253,
		Longitude: -59.523613,
		CCTLD:     `.fk`,
	},
	Micronesia: {
		Name:      `Micronesia`,
		Code:      Micronesia,
		Latitude:  7.425554,
		Longitude: 150.550812,
		CCTLD:     `.fm`,
	},
	FaroeIslands: {
		Name:      `Faroe Islands`,
		Code:      FaroeIslands,
		Latitude:  61.892635,
		Longitude: -6.911806,
		CCTLD:     `.fo`,
	},
	France: {
		Name:      `France`,
		Code:      France,
		Latitude:  46.227638,
		Longitude: 2.213749,
		CCTLD:     `.fr`,
	},
	Gabon: {
		Name:      `Gabon`,
		Code:      Gabon,
		Latitude:  -0.803689,
		Longitude: 11.609444,
		CCTLD:     `.ga`,
	},
	UnitedKingdom: {
		Name:      `United Kingdom`,
		Code:      UnitedKingdom,
		Latitude:  55.378051,
		Longitude: -3.435973,
		CCTLD:     `.uk`,
	},
	Grenada: {
		Name:      `Grenada`,
		Code:      Grenada,
		Latitude:  12.262776,
		Longitude: -61.604171,
		CCTLD:     `.gd`,
	},
	Georgia: {
		Name:      `Georgia`,
		Code:      Georgia,
		Latitude:  42.315407,
		Longitude: 43.356892,
		CCTLD:     `.ge`,
	},
	FrenchGuiana: {
		Name:      `French Guiana`,
		Code:      FrenchGuiana,
		Latitude:  3.933889,
		Longitude: -53.125782,
		CCTLD:     `.gf`,
	},
	Guernsey: {
		Name:      `Guernsey`,
		Code:      Guernsey,
		Latitude:  49.465691,
		Longitude: -2.585278,
		CCTLD:     `.gg`,
	},
	Ghana: {
		Name:      `Ghana`,
		Code:      Ghana,
		Latitude:  7.946527,
		Longitude: -1.023194,
		CCTLD:     `.gh`,
	},
	Gibraltar: {
		Name:      `Gibraltar`,
		Code:      Gibraltar,
		Latitude:  36.137741,
		Longitude: -5.345374,
		CCTLD:     `.gi`,
	},
	Greenland: {
		Name:      `Greenland`,
		Code:      Greenland,
		Latitude:  71.706936,
		Longitude: -42.604303,
		CCTLD:     `.gl`,
	},
	Gambia: {
		Name:      `Gambia`,
		Code:      Gambia,
		Latitude:  13.443182,
		Longitude: -15.310139,
		CCTLD:     `.gm`,
	},
	Guinea: {
		Name:      `Guinea`,
		Code:      Guinea,
		Latitude:  9.945587,
		Longitude: -9.696645,
		CCTLD:     `.gn`,
	},
	Guadeloupe: {
		Name:      `Guadeloupe`,
		Code:      Guadeloupe,
		Latitude:  16.995971,
		Longitude: -62.067641,
		CCTLD:     `.gp`,
	},
	EquatorialGuinea: {
		Name:      `Equatorial Guinea`,
		Code:      EquatorialGuinea,
		Latitude:  1.650801,
		Longitude: 10.267895,
		CCTLD:     `.gq`,
	},
	Greece: {
		Name:      `Greece`,
		Code:      Greece,
		Latitude:  39.074208,
		Longitude: 21.824312,
		CCTLD:     `.gr`,
	},
	SouthGeorgiaAndTheSouthSandwichIslands: {
		Name:      `South Georgia and the South Sandwich Islands`,
		Code:      SouthGeorgiaAndTheSouthSandwichIslands,
		Latitude:  -54.429579,
		Longitude: -36.587909,
		CCTLD:     `.gs`,
	},
	Guatemala: {
		Name:      `Guatemala`,
		Code:      Guatemala,
		Latitude:  15.783471,
		Longitude: -90.230759,
		CCTLD:     `.gt`,
	},
	Guam: {
		Name:      `Guam`,
		Code:      Guam,
		Latitude:  13.444304,
		Longitude: 144.793731,
		CCTLD:     `.gu`,
	},
	GuineaBissau: {
		Name:      `Guinea-Bissau`,
		Code:      GuineaBissau,
		Latitude:  11.803749,
		Longitude: -15.180413,
		CCTLD:     `.gw`,
	},
	Guyana: {
		Name:      `Guyana`,
		Code:      Guyana,
		Latitude:  4.860416,
		Longitude: -58.93018,
		CCTLD:     `.gy`,
	},
	GazaStrip: {
		Name:      `Gaza Strip`,
		Code:      GazaStrip,
		Latitude:  31.354676,
		Longitude: 34.308825,
	},
	HongKong: {
		Name:      `Hong Kong`,
		Code:      HongKong,
		Latitude:  22.396428,
		Longitude: 114.109497,
		CCTLD:     `.hk`,
	},
	HeardIslandAndMcDonaldIslands: {
		Name:      `Heard Island and McDonald Islands`,
		Code:      HeardIslandAndMcDonaldIslands,
		Latitude:  -53.08181,
		Longitude: 73.504158,
		CCTLD:     `.hm`,
	},
	Honduras: {
		Name:      `Honduras`,
		Code:      Honduras,
		Latitude:  15.199999,
		Longitude: -86.241905,
		CCTLD:     `.hn`,
	},
	Croatia: {
		Name:      `Croatia`,
		Code:      Croatia,
		Latitude:  45.1,
		Longitude: 15.2,
		CCTLD:     `.hr`,
	},
	Haiti: {
		Name:      `Haiti`,
		Code:      Haiti,
		Latitude:  18.971187,
		Longitude: -72.285215,
		CCTLD:     `.ht`,
	},
	Hungary: {
		Name:      `Hungary`,
		Code:      Hungary,
		Latitude:  47.162494,
		Longitude: 19.503304,
		CCTLD:     `.hu`,
	},
	Indonesia: {
		Name:      `Indonesia`,
		Code:      Indonesia,
		Latitude:  -0.789275,
		Longitude: 113.921327,
		CCTLD:     `.id`,
	},
	Ireland: {
		Name:      `Ireland`,
		Code:      Ireland,
		Latitude:  53.41291,
		Longitude: -8.24389,
		CCTLD:     `.ie`,
	},
	Israel: {
		Name:      `Israel`,
		Code:      Israel,
		Latitude:  31.046051,
		Longitude: 34.851612,
		CCTLD:     `.il`,
	},
	IsleOfMan: {
		Name:      `Isle of Man`,
		Code:      IsleOfMan,
		Latitude:  54.236107,
		Longitude: -4.548056,
		CCTLD:     `.im`,
	},
	India: {
		Name:      `India`,
		Code:      India,
		Latitude:  20.593684,
		Longitude: 78.96288,
		CCTLD:     `.in`,
	},
	BritishIndianOceanTerritory: {
		Name:      `British Indian Ocean Territory`,
		Code:      BritishIndianOceanTerritory,
		Latitude:  -6.343194,
		Longitude: 71.876519,
		CCTLD:     `.io`,
	},
	Iraq: {
		Name:      `Iraq`,
		Code:      Iraq,
		Latitude:  33.223191,
		Longitude: 43.679291,
		CCTLD:     `.iq`,
	},
	Iran: {
		Name:      `Iran`,
		Code:      Iran,
		Latitude:  32.427908,
		Longitude: 53.688046,
		CCTLD:     `.ir`,
	},
	Iceland: {
		Name:      `Iceland`,
		Code:      Iceland,
		Latitude:  64.963051,
		Longitude: -19.020835,
		CCTLD:     `.is`,
	},
	Italy: {
		Name:      `Italy`,
		Code:      Italy,
		Latitude:  41.87194,
		Longitude: 12.56738,
		CCTLD:     `.it`,
	},
	Jersey: {
		Name:      `Jersey`,
		Code:      Jersey,
		Latitude:  49.214439,
		Longitude: -2.13125,
		CCTLD:     `.je`,
	},
	Jamaica: {
		Name:      `Jamaica`,
		Code:      Jamaica,
		Latitude:  18.109581,
		Longitude: -77.297508,
		CCTLD:     `.jm`,
	},
	Jordan: {
		Name:      `Jordan`,
		Code:      Jordan,
		Latitude:  30.585164,
		Longitude: 36.238414,
		CCTLD:     `.jo`,
	},
	Japan: {
		Name:      `Japan`,
		Code:      Japan,
		Latitude:  36.204824,
		Longitude: 138.252924,
		CCTLD:     `.jp`,
	},
	Kenya: {
		Name:      `Kenya`,
		Code:      Kenya,
		Latitude:  -0.023559,
		Longitude: 37.906193,
		CCTLD:     `.ke`,
	},
	Kyrgyzstan: {
		Name:      `Kyrgyzstan`,
		Code:      Kyrgyzstan,
		Latitude:  41.20438,
		Longitude: 74.766098,
		CCTLD:     `.kg`,
	},
	Cambodia: {
		Name:      `Cambodia`,
		Code:      Cambodia,
		Latitude:  12.565679,
		Longitude: 104.990963,
		CCTLD:     `.kh`,
	},
	Kiribati: {
		Name:      `Kiribati`,
		Code:      Kiribati,
		Latitude:  -3.370417,
		Longitude: -168.734039,
		CCTLD:     `.ki`,
	},
	Comoros: {
		Name:      `Comoros`,
		Code:      Comoros,
		Latitude:  -11.875001,
		Longitude: 43.872219,
		CCTLD:     `.km`,
	},
	SaintKittsAndNevis: {
		Name:      `Saint Kitts and Nevis`,
		Code:      SaintKittsAndNevis,
		Latitude:  17.357822,
		Longitude: -62.782998,
		CCTLD:     `.kn`,
	},
	NorthKorea: {
		Name:      `North Korea`,
		Code:      NorthKorea,
		Latitude:  40.339852,
		Longitude: 127.510093,
		CCTLD:     `.kp`,
	},
	SouthKorea: {
		Name:      `South Korea`,
		Code:      SouthKorea,
		Latitude:  35.907757,
		Longitude: 127.766922,
		CCTLD:     `.kr`,
	},
	Kuwait: {
		Name:      `Kuwait`,
		Code:      Kuwait,
		Latitude:  29.31166,
		Longitude: 47.481766,
		CCTLD:     `.kw`,
	},
	CaymanIslands: {
		Name:      `Cayman Islands`,
		Code:      CaymanIslands,
		Latitude:  19.513469,
		Longitude: -80.566956,
		CCTLD:     `.ky`,
	},
	Kazakhstan: {
		Name:      `Kazakhstan`,
		Code:      Kazakhstan,
		Latitude:  48.019573,
		Longitude: 66.923684,
		CCTLD:     `.kz`,
	},
	Laos: {
		Name:      `Laos`,
		Code:      Laos,
		Latitude:  19.85627,
		Longitude: 102.495496,
		CCTLD:     `.la`,
	},
	Lebanon: {
		Name:      `Lebanon`,
		Code:      Lebanon,
		Latitude:  33.854721,
		Longitude: 35.862285,
		CCTLD:     `.lb`,
	},
	SaintLucia: {
		Name:      `Saint Lucia`,
		Code:      SaintLucia,
		Latitude:  13.909444,
		Longitude: -60.978893,
		CCTLD:     `.lc`,
	},
	Liechtenstein: {
		Name:      `Liechtenstein`,
		Code:      Liechtenstein,
		Latitude:  47.166,
		Longitude: 9.555373,
		CCTLD:     `.li`,
	},
	SriLanka: {
		Name:      `Sri Lanka`,
		Code:      SriLanka,
		Latitude:  7.873054,
		Longitude: 80.771797,
		CCTLD:     `.lk`,
	},
	Liberia: {
		Name:      `Liberia`,
		Code:      Liberia,
		Latitude:  6.428055,
		Longitude: -9.429499,
		CCTLD:     `.lr`,
	},
	Lesotho: {
		Name:      `Lesotho`,
		Code:      Lesotho,
		Latitude:  -29.609988,
		Longitude: 28.233608,
		CCTLD:     `.ls`,
	},
	Lithuania: {
		Name:      `Lithuania`,
		Code:      Lithuania,
		Latitude:  55.169438,
		Longitude: 23.881275,
		CCTLD:     `.lt`,
	},
	Luxembourg: {
		Name:      `Luxembourg`,
		Code:      Luxembourg,
		Latitude:  49.815273,
		Longitude: 6.129583,
		CCTLD:     `.lu`,
	},
	Latvia: {
		Name:      `Latvia`,
		Code:      Latvia,
		Latitude:  56.879635,
		Longitude: 24.603189,
		CCTLD:     `.lv`,
	},
	Libya: {
		Name:      `Libya`,
		Code:      Libya,
		Latitude:  26.3351,
		Longitude: 17.228331,
		CCTLD:     `.ly`,
	},
	Morocco: {
		Name:      `Morocco`,
		Code:      Morocco,
		Latitude:  31.791702,
		Longitude: -7.09262,
		CCTLD:     `.ma`,
	},
	Monaco: {
		Name:      `Monaco`,
		Code:      Monaco,
		Latitude:  43.750298,
		Longitude: 7.412841,
		CCTLD:     `.mc`,
	},
	Moldova: {
		Name:      `Moldova`,
		Code:      Moldova,
		Latitude:  47.411631,
		Longitude: 28.369885,
		CCTLD:     `.md`,
	},
	Montenegro: {
		Name:      `Montenegro`,
		Code:      Montenegro,
		Latitude:  42.708678,
		Longitude: 19.37439,
		CCTLD:     `.me`,
	},
	Madagascar: {
		Name:      `Madagascar`,
		Code:      Madagascar,
		Latitude:  -18.766947,
		Longitude: 46.869107,
		CCTLD:     `.mg`,
	},
	MarshallIslands: {
		Name:      `Marshall Islands`,
		Code:      MarshallIslands,
		Latitude:  7.131474,
		Longitude: 171.184478,
		CCTLD:     `.mh`,
	},
	Macedonia: {
		Name:      `Macedonia`,
		Code:      Macedonia,
		Latitude:  41.608635,
		Longitude: 21.745275,
		CCTLD:     `.mk`,
	},
	Mali: {
		Name:      `Mali`,
		Code:      Mali,
		Latitude:  17.570692,
		Longitude: -3.996166,
		CCTLD:     `.ml`,
	},
	Myanmar: {
		Name:      `Myanmar (Burma)`,
		Code:      Myanmar,
		Latitude:  21.913965,
		Longitude: 95.956223,
		CCTLD:     `.mm`,
	},
	Mongolia: {
		Name:      `Mongolia`,
		Code:      Mongolia,
		Latitude:  46.862496,
		Longitude: 103.846656,
		CCTLD:     `.mn`,
	},
	Macau: {
		Name:      `Macau`,
		Code:      Macau,
		Latitude:  22.198745,
		Longitude: 113.543873,
		CCTLD:     `.mo`,
	},
	NorthernMarianaIslands: {
		Name:      `Northern Mariana Islands`,
		Code:      NorthernMarianaIslands,
		Latitude:  17.33083,
		Longitude: 145.38469,
		CCTLD:     `.mp`,
	},
	Martinique: {
		Name:      `Martinique`,
		Code:      Martinique,
		Latitude:  14.641528,
		Longitude: -61.024174,
		CCTLD:     `.mq`,
	},
	Mauritania: {
		Name:      `Mauritania`,
		Code:      Mauritania,
		Latitude:  21.00789,
		Longitude: -10.940835,
		CCTLD:     `.mr`,
	},
	Montserrat: {
		Name:      `Montserrat`,
		Code:      Montserrat,
		Latitude:  16.742498,
		Longitude: -62.187366,
		CCTLD:     `.ms`,
	},
	Malta: {
		Name:      `Malta`,
		Code:      Malta,
		Latitude:  35.937496,
		Longitude: 14.375416,
		CCTLD:     `.mt`,
	},
	Mauritius: {
		Name:      `Mauritius`,
		Code:      Mauritius,
		Latitude:  -20.348404,
		Longitude: 57.552152,
		CCTLD:     `.mu`,
	},
	Maldives: {
		Name:      `Maldives`,
		Code:      Maldives,
		Latitude:  3.202778,
		Longitude: 73.22068,
		CCTLD:     `.mv`,
	},
	Malawi: {
		Name:      `Malawi`,
		Code:      Malawi,
		Latitude:  -13.254308,
		Longitude: 34.301525,
		CCTLD:     `.mw`,
	},
	Mexico: {
		Name:      `Mexico`,
		Code:      Mexico,
		Latitude:  23.634501,
		Longitude: -102.552784,
		CCTLD:     `.mx`,
	},
	Malaysia: {
		Name:      `Malaysia`,
		Code:      Malaysia,
		Latitude:  4.210484,
		Longitude: 101.975766,
		CCTLD:     `.my`,
	},
	Mozambique: {
		Name:      `Mozambique`,
		Code:      Mozambique,
		Latitude:  -18.665695,
		Longitude: 35.529562,
		CCTLD:     `.mz`,
	},
	Namibia: {
		Name:      `Namibia`,
		Code:      Namibia,
		Latitude:  -22.95764,
		Longitude: 18.49041,
		CCTLD:     `.na`,
	},
	NewCaledonia: {
		Name:      `New Caledonia`,
		Code:      NewCaledonia,
		Latitude:  -20.904305,
		Longitude: 165.618042,
		CCTLD:     `.nc`,
	},
	Niger: {
		Name:      `Niger`,
		Code:      Niger,
		Latitude:  17.607789,
		Longitude: 8.081666,
		CCTLD:     `.ne`,
	},
	NorfolkIsland: {
		Name:      `Norfolk Island`,
		Code:      NorfolkIsland,
		Latitude:  -29.040835,
		Longitude: 167.954712,
		CCTLD:     `.nf`,
	},
	Nigeria: {
		Name:      `Nigeria`,
		Code:      Nigeria,
		Latitude:  9.081999,
		Longitude: 8.675277,
		CCTLD:     `.ng`,
	},
	Nicaragua: {
		Name:      `Nicaragua`,
		Code:      Nicaragua,
		Latitude:  12.865416,
		Longitude: -85.207229,
		CCTLD:     `.ni`,
	},
	Netherlands: {
		Name:      `Netherlands`,
		Code:      Netherlands,
		Latitude:  52.132633,
		Longitude: 5.291266,
		CCTLD:     `.nl`,
	},
	Norway: {
		Name:      `Norway`,
		Code:      Norway,
		Latitude:  60.472024,
		Longitude: 8.468946,
		CCTLD:     `.no`,
	},
	Nepal: {
		Name:      `Nepal`,
		Code:      Nepal,
		Latitude:  28.394857,
		Longitude: 84.124008,
		CCTLD:     `.np`,
	},
	Nauru: {
		Name:      `Nauru`,
		Code:      Nauru,
		Latitude:  -0.522778,
		Longitude: 166.931503,
		CCTLD:     `.nr`,
	},
	Niue: {
		Name:      `Niue`,
		Code:      Niue,
		Latitude:  -19.054445,
		Longitude: -169.867233,
		CCTLD:     `.nu`,
	},
	NewZealand: {
		Name:      `New Zealand`,
		Code:      NewZealand,
		Latitude:  -40.900557,
		Longitude: 174.885971,
		CCTLD:     `.nz`,
	},
	Oman: {
		Name:      `Oman`,
		Code:      Oman,
		Latitude:  21.512583,
		Longitude: 55.923255,
		CCTLD:     `.om`,
	},
	Panama: {
		Name:      `Panama`,
		Code:      Panama,
		Latitude:  8.537981,
		Longitude: -80.782127,
		CCTLD:     `.pa`,
	},
	Peru: {
		Name:      `Peru`,
		Code:      Peru,
		Latitude:  -9.189967,
		Longitude: -75.015152,
		CCTLD:     `.pe`,
	},
	FrenchPolynesia: {
		Name:      `French Polynesia`,
		Code:      FrenchPolynesia,
		Latitude:  -17.679742,
		Longitude: -149.406843,
		CCTLD:     `.pf`,
	},
	PapuaNewGuinea: {
		Name:      `Papua New Guinea`,
		Code:      PapuaNewGuinea,
		Latitude:  -6.314993,
		Longitude: 143.95555,
		CCTLD:     `.pg`,
	},
	Philippines: {
		Name:      `Philippines`,
		Code:      Philippines,
		Latitude:  12.879721,
		Longitude: 121.774017,
		CCTLD:     `.ph`,
	},
	Pakistan: {
		Name:      `Pakistan`,
		Code:      Pakistan,
		Latitude:  30.375321,
		Longitude: 69.345116,
		CCTLD:     `.pk`,
	},
	Poland: {
		Name:      `Poland`,
		Code:      Poland,
		Latitude:  51.919438,
		Longitude: 19.145136,
		CCTLD:     `.pl`,
	},
	SaintPierreAndMiquelon: {
		Name:      `Saint Pierre and Miquelon`,
		Code:      SaintPierreAndMiquelon,
		Latitude:  46.941936,
		Longitude: -56.27111,
		CCTLD:     `.pm`,
	},
	PitcairnIslands: {
		Name:      `Pitcairn Islands`,
		Code:      PitcairnIslands,
		Latitude:  -24.703615,
		Longitude: -127.439308,
		CCTLD:     `.pn`,
	},
	PuertoRico: {
		Name:      `Puerto Rico`,
		Code:      PuertoRico,
		Latitude:  18.220833,
		Longitude: -66.590149,
		CCTLD:     `.pr`,
	},
	Palestine: {
		Name:      `Palestine`,
		Code:      Palestine,
		Latitude:  31.952162,
		Longitude: 35.233154,
		CCTLD:     `.ps`,
	},
	Portugal: {
		Name:      `Portugal`,
		Code:      Portugal,
		Latitude:  39.399872,
		Longitude: -8.224454,
		CCTLD:     `.pt`,
	},
	Palau: {
		Name:      `Palau`,
		Code:      Palau,
		Latitude:  7.51498,
		Longitude: 134.58252,
		CCTLD:     `.pw`,
	},
	Paraguay: {
		Name:      `Paraguay`,
		Code:      Paraguay,
		Latitude:  -23.442503,
		Longitude: -58.443832,
		CCTLD:     `.py`,
	},
	Qatar: {
		Name:      `Qatar`,
		Code:      Qatar,
		Latitude:  25.354826,
		Longitude: 51.183884,
		CCTLD:     `.qa`,
	},
	Reunion: {
		Name:      `Réunion`,
		Code:      Reunion,
		Latitude:  -21.115141,
		Longitude: 55.536384,
		CCTLD:     `.re`,
	},
	Romania: {
		Name:      `Romania`,
		Code:      Romania,
		Latitude:  45.943161,
		Longitude: 24.96676,
		CCTLD:     `.ro`,
	},
	Serbia: {
		Name:      `Serbia`,
		Code:      Serbia,
		Latitude:  44.016521,
		Longitude: 21.005859,
		CCTLD:     `.rs`,
	},
	Russia: {
		Name:      `Russia`,
		Code:      Russia,
		Latitude:  61.52401,
		Longitude: 105.318756,
		CCTLD:     `.ru`,
	},
	Rwanda: {
		Name:      `Rwanda`,
		Code:      Rwanda,
		Latitude:  -1.940278,
		Longitude: 29.873888,
		CCTLD:     `.rw`,
	},
	SaudiArabia: {
		Name:      `Saudi Arabia`,
		Code:      SaudiArabia,
		Latitude:  23.885942,
		Longitude: 45.079162,
		CCTLD:     `.sa`,
	},
	SolomonIslands: {
		Name:      `Solomon Islands`,
		Code:      SolomonIslands,
		Latitude:  -9.64571,
		Longitude: 160.156194,
		CCTLD:     `.sb`,
	},
	Seychelles: {
		Name:      `Seychelles`,
		Code:      Seychelles,
		Latitude:  -4.679574,
		Longitude: 55.491977,
		CCTLD:     `.sc`,
	},
	Sudan: {
		Name:      `Sudan`,
		Code:      Sudan,
		Latitude:  12.862807,
		Longitude: 30.217636,
		CCTLD:     `.sd`,
	},
	Sweden: {
		Name:      `Sweden`,
		Code:      Sweden,
		Latitude:  60.128161,
		Longitude: 18.643501,
		CCTLD:     `.se`,
	},
	Singapore: {
		Name:      `Singapore`,
		Code:      Singapore,
		Latitude:  1.352083,
		Longitude: 103.819836,
		CCTLD:     `.sg`,
	},
	SaintHelena: {
		Name:      `Saint Helena`,
		Code:      SaintHelena,
		Latitude:  -24.143474,
		Longitude: -10.030696,
		CCTLD:     `.sh`,
	},
	Slovenia: {
		Name:      `Slovenia`,
		Code:      Slovenia,
		Latitude:  46.151241,
		Longitude: 14.995463,
		CCTLD:     `.si`,
	},
	SvalbardAndJanMayen: {
		Name:      `Svalbard and Jan Mayen`,
		Code:      SvalbardAndJanMayen,
		Latitude:  77.553604,
		Longitude: 23.670272,
	},
	Slovakia: {
		Name:      `Slovakia`,
		Code:      Slovakia,
		Latitude:  48.669026,
		Longitude: 19.699024,
		CCTLD:     `.sk`,
	},
	SierraLeone: {
		Name:      `Sierra Leone`,
		Code:      SierraLeone,
		Latitude:  8.460555,
		Longitude: -11.779889,
		CCTLD:     `.sl`,
	},
	SanMarino: {
		Name:      `San Marino`,
		Code:      SanMarino,
		Latitude:  43.94236,
		Longitude: 12.457777,
		CCTLD:     `.sm`,
	},
	Senegal: {
		Name:      `Senegal`,
		Code:      Senegal,
		Latitude:  14.497401,
		Longitude: -14.452362,
		CCTLD:     `.sn`,
	},
	Somalia: {
		Name:      `Somalia`,
		Code:      Somalia,
		Latitude:  5.152149,
		Longitude: 46.199616,
		CCTLD:     `.so`,
	},
	Suriname: {
		Name:      `Suriname`,
		Code:      Suriname,
		Latitude:  3.919305,
		Longitude: -56.027783,
		CCTLD:     `.sr`,
	},
	SouthSudan: {
		Name:  `South Sudan`,
		Code:  SouthSudan,
		CCTLD: `.ss`,
	},
	SaoTomeAndPrincipe: {
		Name:      `São Tomé and Príncipe`,
		Code:      SaoTomeAndPrincipe,
		Latitude:  0.18636,
		Longitude: 6.613081,
		CCTLD:     `.st`,
	},
	SovietUnion: {
		Name:  `Soviet Union`,
		Code:  SovietUnion,
		CCTLD: `.su`,
	},
	ElSalvador: {
		Name:      `El Salvador`,
		Code:      ElSalvador,
		Latitude:  13.794185,
		Longitude: -88.89653,
		CCTLD:     `.sv`,
	},
	Syria: {
		Name:      `Syria`,
		Code:      Syria,
		Latitude:  34.802075,
		Longitude: 38.996815,
		CCTLD:     `.sy`,
	},
	Swaziland: {
		Name:      `Swaziland`,
		Code:      Swaziland,
		Latitude:  -26.522503,
		Longitude: 31.465866,
		CCTLD:     `.sz`,
	},
	TurksAndCaicosIslands: {
		Name:      `Turks and Caicos Islands`,
		Code:      TurksAndCaicosIslands,
		Latitude:  21.694025,
		Longitude: -71.797928,
		CCTLD:     `.tc`,
	},
	Chad: {
		Name:      `Chad`,
		Code:      Chad,
		Latitude:  15.454166,
		Longitude: 18.732207,
		CCTLD:     `.td`,
	},
	FrenchSouthernTerritories: {
		Name:      `French Southern Territories`,
		Code:      FrenchSouthernTerritories,
		Latitude:  -49.280366,
		Longitude: 69.348557,
		CCTLD:     `.tf`,
	},
	Togo: {
		Name:      `Togo`,
		Code:      Togo,
		Latitude:  8.619543,
		Longitude: 0.824782,
		CCTLD:     `.tg`,
	},
	Thailand: {
		Name:      `Thailand`,
		Code:      Thailand,
		Latitude:  15.870032,
		Longitude: 100.992541,
		CCTLD:     `.th`,
	},
	Tajikistan: {
		Name:      `Tajikistan`,
		Code:      Tajikistan,
		Latitude:  38.861034,
		Longitude: 71.276093,
		CCTLD:     `.tj`,
	},
	Tokelau: {
		Name:      `Tokelau`,
		Code:      Tokelau,
		Latitude:  -8.967363,
		Longitude: -171.855881,
		CCTLD:     `.tk`,
	},
	TimorLeste: {
		Name:      `Timor-Leste`,
		Code:      TimorLeste,
		Latitude:  -8.874217,
		Longitude: 125.727539,
		CCTLD:     `.tl`,
	},
	Turkmenistan: {
		Name:      `Turkmenistan`,
		Code:      Turkmenistan,
		Latitude:  38.969719,
		Longitude: 59.556278,
		CCTLD:     `.tm`,
	},
	Tunisia: {
		Name:      `Tunisia`,
		Code:      Tunisia,
		Latitude:  33.886917,
		Longitude: 9.537499,
		CCTLD:     `.tn`,
	},
	Tonga: {
		Name:      `Tonga`,
		Code:      Tonga,
		Latitude:  -21.178986,
		Longitude: -175.198242,
		CCTLD:     `.to`,
	},
	Turkey: {
		Name:      `Turkey`,
		Code:      Turkey,
		Latitude:  38.963745,
		Longitude: 35.243322,
		CCTLD:     `.tr`,
	},
	TrinidadAndTobago: {
		Name:      `Trinidad and Tobago`,
		Code:      TrinidadAndTobago,
		Latitude:  10.691803,
		Longitude: -61.222503,
		CCTLD:     `.tt`,
	},
	Tuvalu: {
		Name:      `Tuvalu`,
		Code:      Tuvalu,
		Latitude:  -7.109535,
		Longitude: 177.64933,
		CCTLD:     `.tv`,
	},
	Taiwan: {
		Name:      `Taiwan`,
		Code:      Taiwan,
		Latitude:  23.69781,
		Longitude: 120.960515,
		CCTLD:     `.tw`,
	},
	Tanzania: {
		Name:      `Tanzania`,
		Code:      Tanzania,
		Latitude:  -6.369028,
		Longitude: 34.888822,
		CCTLD:     `.tz`,
	},
	Ukraine: {
		Name:      `Ukraine`,
		Code:      Ukraine,
		Latitude:  48.379433,
		Longitude: 31.16558,
		CCTLD:     `.ua`,
	},
	Uganda: {
		Name:      `Uganda`,
		Code:      Uganda,
		Latitude:  1.373333,
		Longitude: 32.290275,
		CCTLD:     `.ug`,
	},
	UnitedStatesMinorOutlyingIslands: {
		Name: `U.S. Minor Outlying Islands`,
		Code: UnitedStatesMinorOutlyingIslands,
	},
	Uruguay: {
		Name:      `Uruguay`,
		Code:      Uruguay,
		Latitude:  -32.522779,
		Longitude: -55.765835,
		CCTLD:     `.uy`,
	},
	Uzbekistan: {
		Name:      `Uzbekistan`,
		Code:      Uzbekistan,
		Latitude:  41.377491,
		Longitude: 64.585262,
		CCTLD:     `.uz`,
	},
	VaticanCity: {
		Name:      `Vatican City`,
		Code:      VaticanCity,
		Latitude:  41.902916,
		Longitude: 12.453389,
		CCTLD:     `.va`,
	},
	SaintVincentAndTheGrenadines: {
		Name:      `Saint Vincent and the Grenadines`,
		Code:      SaintVincentAndTheGrenadines,
		Latitude:  12.984305,
		Longitude: -61.287228,
		CCTLD:     `.vc`,
	},
	Venezuela: {
		Name:      `Venezuela`,
		Code:      Venezuela,
		Latitude:  6.42375,
		Longitude: -66.58973,
		CCTLD:     `.ve`,
	},
	BritishVirginIslands: {
		Name:      `British Virgin Islands`,
		Code:      BritishVirginIslands,
		Latitude:  18.420695,
		Longitude: -64.639968,
		CCTLD:     `.vg`,
	},
	UnitedStates: {
		Name:      `United States of America`,
		Code:      UnitedStates,
		Latitude:  37.09024,
		Longitude: -95.712891,
		CCTLD:     `.us`,
	},
	UnitedStatesVirginIslands: {
		Name:      `U.S. Virgin Islands`,
		Code:      UnitedStatesVirginIslands,
		Latitude:  18.335765,
		Longitude: -64.896335,
		CCTLD:     `.vi`,
	},
	Vietnam: {
		Name:      `Vietnam`,
		Code:      Vietnam,
		Latitude:  14.058324,
		Longitude: 108.277199,
		CCTLD:     `.vn`,
	},
	Vanuatu: {
		Name:      `Vanuatu`,
		Code:      Vanuatu,
		Latitude:  -15.376706,
		Longitude: 166.959158,
		CCTLD:     `.vu`,
	},
	WallisAndFutuna: {
		Name:      `Wallis and Futuna`,
		Code:      WallisAndFutuna,
		Latitude:  -13.768752,
		Longitude: -177.156097,
		CCTLD:     `.wf`,
	},
	Samoa: {
		Name:      `Samoa`,
		Code:      Samoa,
		Latitude:  -13.759029,
		Longitude: -172.104629,
		CCTLD:     `.ws`,
	},
	Kosovo: {
		Name:      `Kosovo`,
		Code:      Kosovo,
		Latitude:  42.602636,
		Longitude: 20.902977,
	},
	Yemen: {
		Name:      `Yemen`,
		Code:      Yemen,
		Latitude:  15.552727,
		Longitude: 48.516388,
		CCTLD:     `.ye`,
	},
	Mayotte: {
		Name:      `Mayotte`,
		Code:      Mayotte,
		Latitude:  -12.8275,
		Longitude: 45.166244,
		CCTLD:     `.yt`,
	},
	SouthAfrica: {
		Name:      `South Africa`,
		Code:      SouthAfrica,
		Latitude:  -30.559482,
		Longitude: 22.937506,
		CCTLD:     `.za`,
	},
	Zambia: {
		Name:      `Zambia`,
		Code:      Zambia,
		Latitude:  -13.133897,
		Longitude: 27.849332,
		CCTLD:     `.zm`,
	},
	Zimbabwe: {
		Name:      `Zimbabwe`,
		Code:      Zimbabwe,
		Latitude:  -19.015438,
		Longitude: 29.154857,
		CCTLD:     `.zw`,
	},
}
