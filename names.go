package stonecutters

import (
	"fmt"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var NAMountains = []string{
	"Denali",
	"MtLogan",
	"PicodeOrizaba",
	"MtStElias",
	"VolcanPopocatepetl",
	"MtForaker",
	"MtLucania",
	"VolcanIztaccihuatl",
	"KingPeak",
	"MtBona",
	"MtSteele",
	"MtBlackburn",
	"MtSanford",
	"MtWood",
	"MtVancouver",
	"NevadodeToluca",
	"MtFairweather",
	"MtHubbard",
	"MtBear",
	"MtWalsh",
	"MtHunter",
	"VolcanLaMalinche",
	"MtWhitney",
	"UniversityPeak",
	"MtElbert",
	"MtHarvard",
	"MtRainier",
	"BlancaPeak",
	"UncompahgrePeak",
	"McArthurPeak",
	"CrestonePeak",
	"MtLincoln",
	"GraysPeak",
	"MtAntero",
	"CastlePeak",
	"MtEvans",
	"LongsPeak",
	"WhiteMountainPeak",
	"MtWilson",
	"NorthPalisade",
	"NevadodeColima",
	"MtPrinceton",
	"MtWrangell",
	"MtShasta",
	"MaroonPeak",
	"MtSneffels",
	"PikesPeak",
	"MtEolus",
	"MtAugusta",
	"CulebraPeak",
	"SanLuisPeak",
	"MtoftheHolyCross",
	"MtHumphreys",
	"MtOuray",
	"MtStrickland",
	"VermilionPeak",
	"AtnaPeaks",
	"RegalMountain",
	"VolcánTajumulco",
	"MtHayes",
	"MtSilverheels",
	"GannettPeak",
	"MtKaweah",
	"VolcanCofredePerote",
	"GrandTeton",
	"MtCook",
	"MtMorgan",
	"MtGabb",
	"BaldMountain",
	"WestSpanishPeak",
	"MtPowell",
	"HaguesPeak",
	"MtDubois",
	"KingsPeak",
	"TreasureMountain",
	"MtPinchot",
	"MtNatazhat",
	"MtJarvis",
	"VolcánTacaná",
	"MtHerard",
	"SummitPeak",
	"AntoraPeak",
	"HesperusMountain",
	"MtSilverthrone",
	"JacquePeak",
	"WindRiverPeak",
	"MtWaddington",
	"MtMarcusBaker",
	"CloudPeak",
	"WheelerPeak",
	"TwilightPeak",
	"FrancsPeak",
	"SouthRiverPeak",
	"MtRitter",
	"BushnellPeak",
	"TruchasPeak",
	"WheelerPeak",
	"MtDana",
	"SpringGlacierPeak",
	"VolcánAcatenango",
}

func notAscii(r rune) bool {
	return r < 32 || r >= 127
}

func normalizeString(str string) (string, error) {
	t := transform.Chain(norm.NFKD, transform.RemoveFunc(notAscii))
	str, _, err := transform.String(t, str)
	if err != nil {
		return "", err
	}
	return str, nil
}

// NormalizedNaMountains returns a North American Mountain name slice of ASCII strings
func NormalizedNaMountains() []string {
	nm := make([]string, 0)
	for _, s := range NAMountains {
		ns, _ := normalizeString(s)
		if ns != "" {
			nm = append(nm, ns)
		}
	}
	return nm
}

// PrefixedNumerics returns a list of ids with a simple name prefix
// followed by an integer separated by a dash. Naming starts at '1' not '0'.
func PrefixedNumerics(prefix string, max int) []string {
	ids := make([]string, 0)
	for i := 1; i <= max; i++ {
		ids = append(ids, fmt.Sprintf("%s%d", prefix, i))
	}
	return ids
}
