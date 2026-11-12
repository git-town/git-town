// Code generated from dialects_builtin.go.jq (make dialects_builtin.go); DO NOT EDIT.

package gherkin

import messages "github.com/cucumber/messages/go/v34"

// Builtin dialects for af (Afrikaans), am (Armenian), an (Aragonese), ar (Arabic), ast (Asturian), az (Azerbaijani), be (Belarusian), bg (Bulgarian), bm (Malay), bs (Bosnian), ca (Catalan), cs (Czech), cy-GB (Welsh), da (Danish), de (German), el (Greek), em (Emoji), en (English), en-Scouse (Scouse), en-au (Australian), en-lol (LOLCAT), en-old (Old English), en-pirate (Pirate), en-tx (Texas), eo (Esperanto), es (Spanish), et (Estonian), fa (Persian), fi (Finnish), fr (French), ga (Irish), gj (Gujarati), gl (Galician), he (Hebrew), hi (Hindi), hr (Croatian), ht (Creole), hu (Hungarian), id (Indonesian), is (Icelandic), it (Italian), ja (Japanese), jv (Javanese), ka (Georgian), kn (Kannada), ko (Korean), lt (Lithuanian), lu (Luxemburgish), lv (Latvian), mk-Cyrl (Macedonian), mk-Latn (Macedonian (Latin)), mn (Mongolian), ne (Nepali), nl (Dutch), no (Norwegian), pa (Panjabi), pl (Polish), pt (Portuguese), ro (Romanian), ru (Russian), sk (Slovak), sl (Slovenian), sr-Cyrl (Serbian), sr-Latn (Serbian (Latin)), sv (Swedish), ta (Tamil), th (Thai), te (Telugu), tlh (Klingon), tr (Turkish), tt (Tatar), uk (Ukrainian), ur (Urdu), uz (Uzbek), vi (Vietnamese), zh-CN (Chinese simplified), ml (Malayalam), zh-TW (Chinese traditional), mr (Marathi), amh (Amharic)
func DialectsBuiltin() DialectProvider {
	return builtinDialects
}

const (
	feature         = "feature"
	rule            = "rule"
	background      = "background"
	scenario        = "scenario"
	scenarioOutline = "scenarioOutline"
	examples        = "examples"
	given           = "given"
	when            = "when"
	then            = "then"
	and             = "and"
	but             = "but"
)

var builtinDialects = gherkinDialectMap{
	"af": &Dialect{
		Language: "af",
		Name:     "Afrikaans",
		Native:   "Afrikaans",
		Keywords: map[string][]string{
			feature: {
				"Funksie",
				"Besigheid Behoefte",
				"Vermoë",
			},
			rule: {
				"Reël",
				"Reel",
			},
			background: {
				"Agtergrond",
			},
			scenario: {
				"Voorbeeld",
				"Situasie",
			},
			scenarioOutline: {
				"Situasie Uiteensetting",
			},
			examples: {
				"Voorbeelde",
			},
			given: {
				"* ",
				"Gegewe ",
			},
			when: {
				"* ",
				"Wanneer ",
			},
			then: {
				"* ",
				"Dan ",
			},
			and: {
				"* ",
				"En ",
			},
			but: {
				"* ",
				"Maar ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Gegewe ": messages.StepKeywordType_CONTEXT,

			"Wanneer ": messages.StepKeywordType_ACTION,

			"Dan ": messages.StepKeywordType_OUTCOME,

			"En ": messages.StepKeywordType_CONJUNCTION,

			"Maar ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"am": &Dialect{
		Language: "am",
		Name:     "Armenian",
		Native:   "հայերեն",
		Keywords: map[string][]string{
			feature: {
				"Ֆունկցիոնալություն",
				"Հատկություն",
			},
			rule: {
				"Rule",
			},
			background: {
				"Կոնտեքստ",
			},
			scenario: {
				"Օրինակ",
				"Սցենար",
			},
			scenarioOutline: {
				"Սցենարի կառուցվացքը",
			},
			examples: {
				"Օրինակներ",
			},
			given: {
				"* ",
				"Դիցուք ",
			},
			when: {
				"* ",
				"Եթե ",
				"Երբ ",
			},
			then: {
				"* ",
				"Ապա ",
			},
			and: {
				"* ",
				"Եվ ",
			},
			but: {
				"* ",
				"Բայց ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Դիցուք ": messages.StepKeywordType_CONTEXT,

			"Եթե ": messages.StepKeywordType_ACTION,

			"Երբ ": messages.StepKeywordType_ACTION,

			"Ապա ": messages.StepKeywordType_OUTCOME,

			"Եվ ": messages.StepKeywordType_CONJUNCTION,

			"Բայց ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"an": &Dialect{
		Language: "an",
		Name:     "Aragonese",
		Native:   "Aragonés",
		Keywords: map[string][]string{
			feature: {
				"Caracteristica",
			},
			rule: {
				"Rule",
			},
			background: {
				"Antecedents",
			},
			scenario: {
				"Eixemplo",
				"Caso",
			},
			scenarioOutline: {
				"Esquema del caso",
			},
			examples: {
				"Eixemplos",
			},
			given: {
				"* ",
				"Dau ",
				"Dada ",
				"Daus ",
				"Dadas ",
			},
			when: {
				"* ",
				"Cuan ",
			},
			then: {
				"* ",
				"Alavez ",
				"Allora ",
				"Antonces ",
			},
			and: {
				"* ",
				"Y ",
				"E ",
			},
			but: {
				"* ",
				"Pero ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dau ": messages.StepKeywordType_CONTEXT,

			"Dada ": messages.StepKeywordType_CONTEXT,

			"Daus ": messages.StepKeywordType_CONTEXT,

			"Dadas ": messages.StepKeywordType_CONTEXT,

			"Cuan ": messages.StepKeywordType_ACTION,

			"Alavez ": messages.StepKeywordType_OUTCOME,

			"Allora ": messages.StepKeywordType_OUTCOME,

			"Antonces ": messages.StepKeywordType_OUTCOME,

			"Y ": messages.StepKeywordType_CONJUNCTION,

			"E ": messages.StepKeywordType_CONJUNCTION,

			"Pero ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ar": &Dialect{
		Language: "ar",
		Name:     "Arabic",
		Native:   "العربية",
		Keywords: map[string][]string{
			feature: {
				"خاصية",
			},
			rule: {
				"Rule",
			},
			background: {
				"الخلفية",
			},
			scenario: {
				"مثال",
				"سيناريو",
			},
			scenarioOutline: {
				"سيناريو مخطط",
			},
			examples: {
				"امثلة",
			},
			given: {
				"* ",
				"بفرض ",
			},
			when: {
				"* ",
				"متى ",
				"عندما ",
			},
			then: {
				"* ",
				"اذاً ",
				"ثم ",
			},
			and: {
				"* ",
				"و ",
			},
			but: {
				"* ",
				"لكن ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"بفرض ": messages.StepKeywordType_CONTEXT,

			"متى ": messages.StepKeywordType_ACTION,

			"عندما ": messages.StepKeywordType_ACTION,

			"اذاً ": messages.StepKeywordType_OUTCOME,

			"ثم ": messages.StepKeywordType_OUTCOME,

			"و ": messages.StepKeywordType_CONJUNCTION,

			"لكن ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ast": &Dialect{
		Language: "ast",
		Name:     "Asturian",
		Native:   "asturianu",
		Keywords: map[string][]string{
			feature: {
				"Carauterística",
			},
			rule: {
				"Rule",
			},
			background: {
				"Antecedentes",
			},
			scenario: {
				"Exemplo",
				"Casu",
			},
			scenarioOutline: {
				"Esbozu del casu",
			},
			examples: {
				"Exemplos",
			},
			given: {
				"* ",
				"Dáu ",
				"Dada ",
				"Daos ",
				"Daes ",
			},
			when: {
				"* ",
				"Cuando ",
			},
			then: {
				"* ",
				"Entós ",
			},
			and: {
				"* ",
				"Y ",
				"Ya ",
			},
			but: {
				"* ",
				"Peru ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dáu ": messages.StepKeywordType_CONTEXT,

			"Dada ": messages.StepKeywordType_CONTEXT,

			"Daos ": messages.StepKeywordType_CONTEXT,

			"Daes ": messages.StepKeywordType_CONTEXT,

			"Cuando ": messages.StepKeywordType_ACTION,

			"Entós ": messages.StepKeywordType_OUTCOME,

			"Y ": messages.StepKeywordType_CONJUNCTION,

			"Ya ": messages.StepKeywordType_CONJUNCTION,

			"Peru ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"az": &Dialect{
		Language: "az",
		Name:     "Azerbaijani",
		Native:   "Azərbaycanca",
		Keywords: map[string][]string{
			feature: {
				"Özəllik",
			},
			rule: {
				"Rule",
			},
			background: {
				"Keçmiş",
				"Kontekst",
			},
			scenario: {
				"Nümunə",
				"Ssenari",
			},
			scenarioOutline: {
				"Ssenarinin strukturu",
			},
			examples: {
				"Nümunələr",
			},
			given: {
				"* ",
				"Tutaq ki ",
				"Verilir ",
			},
			when: {
				"* ",
				"Əgər ",
				"Nə vaxt ki ",
			},
			then: {
				"* ",
				"O halda ",
			},
			and: {
				"* ",
				"Və ",
				"Həm ",
			},
			but: {
				"* ",
				"Amma ",
				"Ancaq ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Tutaq ki ": messages.StepKeywordType_CONTEXT,

			"Verilir ": messages.StepKeywordType_CONTEXT,

			"Əgər ": messages.StepKeywordType_ACTION,

			"Nə vaxt ki ": messages.StepKeywordType_ACTION,

			"O halda ": messages.StepKeywordType_OUTCOME,

			"Və ": messages.StepKeywordType_CONJUNCTION,

			"Həm ": messages.StepKeywordType_CONJUNCTION,

			"Amma ": messages.StepKeywordType_CONJUNCTION,

			"Ancaq ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"be": &Dialect{
		Language: "be",
		Name:     "Belarusian",
		Native:   "Беларуская",
		Keywords: map[string][]string{
			feature: {
				"Функцыянальнасць",
				"Фіча",
			},
			rule: {
				"Правілы",
			},
			background: {
				"Кантэкст",
			},
			scenario: {
				"Сцэнарый",
				"Cцэнар",
			},
			scenarioOutline: {
				"Шаблон сцэнарыя",
				"Узор сцэнара",
			},
			examples: {
				"Прыклады",
			},
			given: {
				"* ",
				"Няхай ",
				"Дадзена ",
			},
			when: {
				"* ",
				"Калі ",
			},
			then: {
				"* ",
				"Тады ",
			},
			and: {
				"* ",
				"I ",
				"Ды ",
				"Таксама ",
			},
			but: {
				"* ",
				"Але ",
				"Інакш ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Няхай ": messages.StepKeywordType_CONTEXT,

			"Дадзена ": messages.StepKeywordType_CONTEXT,

			"Калі ": messages.StepKeywordType_ACTION,

			"Тады ": messages.StepKeywordType_OUTCOME,

			"I ": messages.StepKeywordType_CONJUNCTION,

			"Ды ": messages.StepKeywordType_CONJUNCTION,

			"Таксама ": messages.StepKeywordType_CONJUNCTION,

			"Але ": messages.StepKeywordType_CONJUNCTION,

			"Інакш ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"bg": &Dialect{
		Language: "bg",
		Name:     "Bulgarian",
		Native:   "български",
		Keywords: map[string][]string{
			feature: {
				"Функционалност",
			},
			rule: {
				"Правило",
			},
			background: {
				"Предистория",
			},
			scenario: {
				"Пример",
				"Сценарий",
			},
			scenarioOutline: {
				"Рамка на сценарий",
			},
			examples: {
				"Примери",
			},
			given: {
				"* ",
				"Дадено ",
			},
			when: {
				"* ",
				"Когато ",
			},
			then: {
				"* ",
				"То ",
			},
			and: {
				"* ",
				"И ",
			},
			but: {
				"* ",
				"Но ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Дадено ": messages.StepKeywordType_CONTEXT,

			"Когато ": messages.StepKeywordType_ACTION,

			"То ": messages.StepKeywordType_OUTCOME,

			"И ": messages.StepKeywordType_CONJUNCTION,

			"Но ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"bm": &Dialect{
		Language: "bm",
		Name:     "Malay",
		Native:   "Bahasa Melayu",
		Keywords: map[string][]string{
			feature: {
				"Fungsi",
			},
			rule: {
				"Rule",
			},
			background: {
				"Latar Belakang",
			},
			scenario: {
				"Senario",
				"Situasi",
				"Keadaan",
			},
			scenarioOutline: {
				"Kerangka Senario",
				"Kerangka Situasi",
				"Kerangka Keadaan",
				"Garis Panduan Senario",
			},
			examples: {
				"Contoh",
			},
			given: {
				"* ",
				"Diberi ",
				"Bagi ",
			},
			when: {
				"* ",
				"Apabila ",
			},
			then: {
				"* ",
				"Maka ",
				"Kemudian ",
			},
			and: {
				"* ",
				"Dan ",
			},
			but: {
				"* ",
				"Tetapi ",
				"Tapi ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Diberi ": messages.StepKeywordType_CONTEXT,

			"Bagi ": messages.StepKeywordType_CONTEXT,

			"Apabila ": messages.StepKeywordType_ACTION,

			"Maka ": messages.StepKeywordType_OUTCOME,

			"Kemudian ": messages.StepKeywordType_OUTCOME,

			"Dan ": messages.StepKeywordType_CONJUNCTION,

			"Tetapi ": messages.StepKeywordType_CONJUNCTION,

			"Tapi ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"bs": &Dialect{
		Language: "bs",
		Name:     "Bosnian",
		Native:   "Bosanski",
		Keywords: map[string][]string{
			feature: {
				"Karakteristika",
			},
			rule: {
				"Rule",
			},
			background: {
				"Pozadina",
			},
			scenario: {
				"Primjer",
				"Scenariju",
				"Scenario",
			},
			scenarioOutline: {
				"Scenariju-obris",
				"Scenario-outline",
			},
			examples: {
				"Primjeri",
			},
			given: {
				"* ",
				"Dato ",
			},
			when: {
				"* ",
				"Kada ",
			},
			then: {
				"* ",
				"Zatim ",
			},
			and: {
				"* ",
				"I ",
				"A ",
			},
			but: {
				"* ",
				"Ali ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dato ": messages.StepKeywordType_CONTEXT,

			"Kada ": messages.StepKeywordType_ACTION,

			"Zatim ": messages.StepKeywordType_OUTCOME,

			"I ": messages.StepKeywordType_CONJUNCTION,

			"A ": messages.StepKeywordType_CONJUNCTION,

			"Ali ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ca": &Dialect{
		Language: "ca",
		Name:     "Catalan",
		Native:   "català",
		Keywords: map[string][]string{
			feature: {
				"Característica",
				"Funcionalitat",
			},
			rule: {
				"Rule",
			},
			background: {
				"Rerefons",
				"Antecedents",
			},
			scenario: {
				"Exemple",
				"Escenari",
			},
			scenarioOutline: {
				"Esquema de l'escenari",
			},
			examples: {
				"Exemples",
			},
			given: {
				"* ",
				"Donat ",
				"Donada ",
				"Atès ",
				"Atesa ",
			},
			when: {
				"* ",
				"Quan ",
			},
			then: {
				"* ",
				"Aleshores ",
				"Cal ",
			},
			and: {
				"* ",
				"I ",
			},
			but: {
				"* ",
				"Però ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Donat ": messages.StepKeywordType_CONTEXT,

			"Donada ": messages.StepKeywordType_CONTEXT,

			"Atès ": messages.StepKeywordType_CONTEXT,

			"Atesa ": messages.StepKeywordType_CONTEXT,

			"Quan ": messages.StepKeywordType_ACTION,

			"Aleshores ": messages.StepKeywordType_OUTCOME,

			"Cal ": messages.StepKeywordType_OUTCOME,

			"I ": messages.StepKeywordType_CONJUNCTION,

			"Però ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"cs": &Dialect{
		Language: "cs",
		Name:     "Czech",
		Native:   "Česky",
		Keywords: map[string][]string{
			feature: {
				"Požadavek",
			},
			rule: {
				"Pravidlo",
			},
			background: {
				"Pozadí",
				"Kontext",
			},
			scenario: {
				"Příklad",
				"Scénář",
			},
			scenarioOutline: {
				"Náčrt Scénáře",
				"Osnova scénáře",
			},
			examples: {
				"Příklady",
			},
			given: {
				"* ",
				"Pokud ",
				"Za předpokladu ",
			},
			when: {
				"* ",
				"Když ",
			},
			then: {
				"* ",
				"Pak ",
			},
			and: {
				"* ",
				"A také ",
				"A ",
			},
			but: {
				"* ",
				"Ale ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Pokud ": messages.StepKeywordType_CONTEXT,

			"Za předpokladu ": messages.StepKeywordType_CONTEXT,

			"Když ": messages.StepKeywordType_ACTION,

			"Pak ": messages.StepKeywordType_OUTCOME,

			"A také ": messages.StepKeywordType_CONJUNCTION,

			"A ": messages.StepKeywordType_CONJUNCTION,

			"Ale ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"cy-GB": &Dialect{
		Language: "cy-GB",
		Name:     "Welsh",
		Native:   "Cymraeg",
		Keywords: map[string][]string{
			feature: {
				"Arwedd",
			},
			rule: {
				"Rule",
			},
			background: {
				"Cefndir",
			},
			scenario: {
				"Enghraifft",
				"Scenario",
			},
			scenarioOutline: {
				"Scenario Amlinellol",
			},
			examples: {
				"Enghreifftiau",
			},
			given: {
				"* ",
				"Anrhegedig a ",
			},
			when: {
				"* ",
				"Pryd ",
			},
			then: {
				"* ",
				"Yna ",
			},
			and: {
				"* ",
				"A ",
			},
			but: {
				"* ",
				"Ond ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Anrhegedig a ": messages.StepKeywordType_CONTEXT,

			"Pryd ": messages.StepKeywordType_ACTION,

			"Yna ": messages.StepKeywordType_OUTCOME,

			"A ": messages.StepKeywordType_CONJUNCTION,

			"Ond ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"da": &Dialect{
		Language: "da",
		Name:     "Danish",
		Native:   "dansk",
		Keywords: map[string][]string{
			feature: {
				"Egenskab",
			},
			rule: {
				"Regel",
			},
			background: {
				"Baggrund",
			},
			scenario: {
				"Eksempel",
				"Scenarie",
			},
			scenarioOutline: {
				"Abstrakt Scenario",
			},
			examples: {
				"Eksempler",
			},
			given: {
				"* ",
				"Givet ",
			},
			when: {
				"* ",
				"Når ",
			},
			then: {
				"* ",
				"Så ",
			},
			and: {
				"* ",
				"Og ",
			},
			but: {
				"* ",
				"Men ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Givet ": messages.StepKeywordType_CONTEXT,

			"Når ": messages.StepKeywordType_ACTION,

			"Så ": messages.StepKeywordType_OUTCOME,

			"Og ": messages.StepKeywordType_CONJUNCTION,

			"Men ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"de": &Dialect{
		Language: "de",
		Name:     "German",
		Native:   "Deutsch",
		Keywords: map[string][]string{
			feature: {
				"Funktionalität",
				"Funktion",
			},
			rule: {
				"Rule",
				"Regel",
			},
			background: {
				"Grundlage",
				"Hintergrund",
				"Voraussetzungen",
				"Vorbedingungen",
			},
			scenario: {
				"Beispiel",
				"Szenario",
			},
			scenarioOutline: {
				"Szenariogrundriss",
				"Szenarien",
			},
			examples: {
				"Beispiele",
			},
			given: {
				"* ",
				"Angenommen ",
				"Gegeben sei ",
				"Gegeben seien ",
			},
			when: {
				"* ",
				"Wenn ",
			},
			then: {
				"* ",
				"Dann ",
			},
			and: {
				"* ",
				"Und ",
			},
			but: {
				"* ",
				"Aber ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Angenommen ": messages.StepKeywordType_CONTEXT,

			"Gegeben sei ": messages.StepKeywordType_CONTEXT,

			"Gegeben seien ": messages.StepKeywordType_CONTEXT,

			"Wenn ": messages.StepKeywordType_ACTION,

			"Dann ": messages.StepKeywordType_OUTCOME,

			"Und ": messages.StepKeywordType_CONJUNCTION,

			"Aber ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"el": &Dialect{
		Language: "el",
		Name:     "Greek",
		Native:   "Ελληνικά",
		Keywords: map[string][]string{
			feature: {
				"Δυνατότητα",
				"Λειτουργία",
			},
			rule: {
				"Rule",
			},
			background: {
				"Υπόβαθρο",
			},
			scenario: {
				"Παράδειγμα",
				"Σενάριο",
			},
			scenarioOutline: {
				"Περιγραφή Σεναρίου",
				"Περίγραμμα Σεναρίου",
			},
			examples: {
				"Παραδείγματα",
				"Σενάρια",
			},
			given: {
				"* ",
				"Δεδομένου ",
			},
			when: {
				"* ",
				"Όταν ",
			},
			then: {
				"* ",
				"Τότε ",
			},
			and: {
				"* ",
				"Και ",
			},
			but: {
				"* ",
				"Αλλά ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Δεδομένου ": messages.StepKeywordType_CONTEXT,

			"Όταν ": messages.StepKeywordType_ACTION,

			"Τότε ": messages.StepKeywordType_OUTCOME,

			"Και ": messages.StepKeywordType_CONJUNCTION,

			"Αλλά ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"em": &Dialect{
		Language: "em",
		Name:     "Emoji",
		Native:   "😀",
		Keywords: map[string][]string{
			feature: {
				"📚",
			},
			rule: {
				"Rule",
			},
			background: {
				"💤",
			},
			scenario: {
				"🥒",
				"📕",
			},
			scenarioOutline: {
				"📖",
			},
			examples: {
				"📓",
			},
			given: {
				"* ",
				"😐",
			},
			when: {
				"* ",
				"🎬",
			},
			then: {
				"* ",
				"🙏",
			},
			and: {
				"* ",
				"😂",
			},
			but: {
				"* ",
				"😔",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"😐": messages.StepKeywordType_CONTEXT,

			"🎬": messages.StepKeywordType_ACTION,

			"🙏": messages.StepKeywordType_OUTCOME,

			"😂": messages.StepKeywordType_CONJUNCTION,

			"😔": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"en": &Dialect{
		Language: "en",
		Name:     "English",
		Native:   "English",
		Keywords: map[string][]string{
			feature: {
				"Feature",
				"Business Need",
				"Ability",
			},
			rule: {
				"Rule",
			},
			background: {
				"Background",
			},
			scenario: {
				"Example",
				"Scenario",
			},
			scenarioOutline: {
				"Scenario Outline",
				"Scenario Template",
			},
			examples: {
				"Examples",
				"Scenarios",
			},
			given: {
				"* ",
				"Given ",
			},
			when: {
				"* ",
				"When ",
			},
			then: {
				"* ",
				"Then ",
			},
			and: {
				"* ",
				"And ",
			},
			but: {
				"* ",
				"But ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Given ": messages.StepKeywordType_CONTEXT,

			"When ": messages.StepKeywordType_ACTION,

			"Then ": messages.StepKeywordType_OUTCOME,

			"And ": messages.StepKeywordType_CONJUNCTION,

			"But ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"en-Scouse": &Dialect{
		Language: "en-Scouse",
		Name:     "Scouse",
		Native:   "Scouse",
		Keywords: map[string][]string{
			feature: {
				"Feature",
			},
			rule: {
				"Rule",
			},
			background: {
				"Dis is what went down",
			},
			scenario: {
				"The thing of it is",
			},
			scenarioOutline: {
				"Wharrimean is",
			},
			examples: {
				"Examples",
			},
			given: {
				"* ",
				"Givun ",
				"Youse know when youse got ",
			},
			when: {
				"* ",
				"Wun ",
				"Youse know like when ",
			},
			then: {
				"* ",
				"Dun ",
				"Den youse gotta ",
			},
			and: {
				"* ",
				"An ",
			},
			but: {
				"* ",
				"Buh ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Givun ": messages.StepKeywordType_CONTEXT,

			"Youse know when youse got ": messages.StepKeywordType_CONTEXT,

			"Wun ": messages.StepKeywordType_ACTION,

			"Youse know like when ": messages.StepKeywordType_ACTION,

			"Dun ": messages.StepKeywordType_OUTCOME,

			"Den youse gotta ": messages.StepKeywordType_OUTCOME,

			"An ": messages.StepKeywordType_CONJUNCTION,

			"Buh ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"en-au": &Dialect{
		Language: "en-au",
		Name:     "Australian",
		Native:   "Australian",
		Keywords: map[string][]string{
			feature: {
				"Pretty much",
			},
			rule: {
				"Rule",
			},
			background: {
				"First off",
			},
			scenario: {
				"Awww, look mate",
			},
			scenarioOutline: {
				"Reckon it's like",
			},
			examples: {
				"You'll wanna",
			},
			given: {
				"* ",
				"Y'know ",
			},
			when: {
				"* ",
				"It's just unbelievable ",
			},
			then: {
				"* ",
				"But at the end of the day I reckon ",
			},
			and: {
				"* ",
				"Too right ",
			},
			but: {
				"* ",
				"Yeah nah ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Y'know ": messages.StepKeywordType_CONTEXT,

			"It's just unbelievable ": messages.StepKeywordType_ACTION,

			"But at the end of the day I reckon ": messages.StepKeywordType_OUTCOME,

			"Too right ": messages.StepKeywordType_CONJUNCTION,

			"Yeah nah ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"en-lol": &Dialect{
		Language: "en-lol",
		Name:     "LOLCAT",
		Native:   "LOLCAT",
		Keywords: map[string][]string{
			feature: {
				"OH HAI",
			},
			rule: {
				"Rule",
			},
			background: {
				"B4",
			},
			scenario: {
				"MISHUN",
			},
			scenarioOutline: {
				"MISHUN SRSLY",
			},
			examples: {
				"EXAMPLZ",
			},
			given: {
				"* ",
				"I CAN HAZ ",
			},
			when: {
				"* ",
				"WEN ",
			},
			then: {
				"* ",
				"DEN ",
			},
			and: {
				"* ",
				"AN ",
			},
			but: {
				"* ",
				"BUT ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"I CAN HAZ ": messages.StepKeywordType_CONTEXT,

			"WEN ": messages.StepKeywordType_ACTION,

			"DEN ": messages.StepKeywordType_OUTCOME,

			"AN ": messages.StepKeywordType_CONJUNCTION,

			"BUT ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"en-old": &Dialect{
		Language: "en-old",
		Name:     "Old English",
		Native:   "Englisc",
		Keywords: map[string][]string{
			feature: {
				"Hwaet",
				"Hwæt",
			},
			rule: {
				"Rule",
			},
			background: {
				"Aer",
				"Ær",
			},
			scenario: {
				"Swa",
			},
			scenarioOutline: {
				"Swa hwaer swa",
				"Swa hwær swa",
			},
			examples: {
				"Se the",
				"Se þe",
				"Se ðe",
			},
			given: {
				"* ",
				"Thurh ",
				"Þurh ",
				"Ðurh ",
			},
			when: {
				"* ",
				"Bæþsealf ",
				"Bæþsealfa ",
				"Bæþsealfe ",
				"Ciricæw ",
				"Ciricæwe ",
				"Ciricæwa ",
			},
			then: {
				"* ",
				"Tha ",
				"Þa ",
				"Ða ",
				"Tha the ",
				"Þa þe ",
				"Ða ðe ",
			},
			and: {
				"* ",
				"Ond ",
				"7 ",
			},
			but: {
				"* ",
				"Ac ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Thurh ": messages.StepKeywordType_CONTEXT,

			"Þurh ": messages.StepKeywordType_CONTEXT,

			"Ðurh ": messages.StepKeywordType_CONTEXT,

			"Bæþsealf ": messages.StepKeywordType_ACTION,

			"Bæþsealfa ": messages.StepKeywordType_ACTION,

			"Bæþsealfe ": messages.StepKeywordType_ACTION,

			"Ciricæw ": messages.StepKeywordType_ACTION,

			"Ciricæwe ": messages.StepKeywordType_ACTION,

			"Ciricæwa ": messages.StepKeywordType_ACTION,

			"Tha ": messages.StepKeywordType_OUTCOME,

			"Þa ": messages.StepKeywordType_OUTCOME,

			"Ða ": messages.StepKeywordType_OUTCOME,

			"Tha the ": messages.StepKeywordType_OUTCOME,

			"Þa þe ": messages.StepKeywordType_OUTCOME,

			"Ða ðe ": messages.StepKeywordType_OUTCOME,

			"Ond ": messages.StepKeywordType_CONJUNCTION,

			"7 ": messages.StepKeywordType_CONJUNCTION,

			"Ac ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"en-pirate": &Dialect{
		Language: "en-pirate",
		Name:     "Pirate",
		Native:   "Pirate",
		Keywords: map[string][]string{
			feature: {
				"Ahoy matey!",
			},
			rule: {
				"Rule",
			},
			background: {
				"Yo-ho-ho",
			},
			scenario: {
				"Heave to",
			},
			scenarioOutline: {
				"Shiver me timbers",
			},
			examples: {
				"Dead men tell no tales",
			},
			given: {
				"* ",
				"Gangway! ",
			},
			when: {
				"* ",
				"Blimey! ",
			},
			then: {
				"* ",
				"Let go and haul ",
			},
			and: {
				"* ",
				"Aye ",
			},
			but: {
				"* ",
				"Avast! ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Gangway! ": messages.StepKeywordType_CONTEXT,

			"Blimey! ": messages.StepKeywordType_ACTION,

			"Let go and haul ": messages.StepKeywordType_OUTCOME,

			"Aye ": messages.StepKeywordType_CONJUNCTION,

			"Avast! ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"en-tx": &Dialect{
		Language: "en-tx",
		Name:     "Texas",
		Native:   "Texas",
		Keywords: map[string][]string{
			feature: {
				"This ain’t my first rodeo",
				"All gussied up",
			},
			rule: {
				"Rule ",
			},
			background: {
				"Lemme tell y'all a story",
			},
			scenario: {
				"All hat and no cattle",
			},
			scenarioOutline: {
				"Serious as a snake bite",
				"Busy as a hound in flea season",
			},
			examples: {
				"Now that's a story longer than a cattle drive in July",
			},
			given: {
				"Fixin' to ",
				"All git out ",
			},
			when: {
				"Quick out of the chute ",
			},
			then: {
				"There’s no tree but bears some fruit ",
			},
			and: {
				"Come hell or high water ",
			},
			but: {
				"Well now hold on, I'll you what ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Fixin' to ": messages.StepKeywordType_CONTEXT,

			"All git out ": messages.StepKeywordType_CONTEXT,

			"Quick out of the chute ": messages.StepKeywordType_ACTION,

			"There’s no tree but bears some fruit ": messages.StepKeywordType_OUTCOME,

			"Come hell or high water ": messages.StepKeywordType_CONJUNCTION,

			"Well now hold on, I'll you what ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"eo": &Dialect{
		Language: "eo",
		Name:     "Esperanto",
		Native:   "Esperanto",
		Keywords: map[string][]string{
			feature: {
				"Trajto",
			},
			rule: {
				"Regulo",
			},
			background: {
				"Fono",
			},
			scenario: {
				"Ekzemplo",
				"Scenaro",
				"Kazo",
			},
			scenarioOutline: {
				"Konturo de la scenaro",
				"Skizo",
				"Kazo-skizo",
			},
			examples: {
				"Ekzemploj",
			},
			given: {
				"* ",
				"Donitaĵo ",
				"Komence ",
			},
			when: {
				"* ",
				"Se ",
			},
			then: {
				"* ",
				"Do ",
			},
			and: {
				"* ",
				"Kaj ",
			},
			but: {
				"* ",
				"Sed ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Donitaĵo ": messages.StepKeywordType_CONTEXT,

			"Komence ": messages.StepKeywordType_CONTEXT,

			"Se ": messages.StepKeywordType_ACTION,

			"Do ": messages.StepKeywordType_OUTCOME,

			"Kaj ": messages.StepKeywordType_CONJUNCTION,

			"Sed ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"es": &Dialect{
		Language: "es",
		Name:     "Spanish",
		Native:   "español",
		Keywords: map[string][]string{
			feature: {
				"Característica",
				"Necesidad del negocio",
				"Requisito",
			},
			rule: {
				"Regla",
				"Regla de negocio",
			},
			background: {
				"Antecedentes",
			},
			scenario: {
				"Ejemplo",
				"Escenario",
			},
			scenarioOutline: {
				"Esquema del escenario",
			},
			examples: {
				"Ejemplos",
			},
			given: {
				"* ",
				"Dado ",
				"Dada ",
				"Dados ",
				"Dadas ",
			},
			when: {
				"* ",
				"Cuando ",
			},
			then: {
				"* ",
				"Entonces ",
			},
			and: {
				"* ",
				"Y ",
				"E ",
			},
			but: {
				"* ",
				"Pero ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dado ": messages.StepKeywordType_CONTEXT,

			"Dada ": messages.StepKeywordType_CONTEXT,

			"Dados ": messages.StepKeywordType_CONTEXT,

			"Dadas ": messages.StepKeywordType_CONTEXT,

			"Cuando ": messages.StepKeywordType_ACTION,

			"Entonces ": messages.StepKeywordType_OUTCOME,

			"Y ": messages.StepKeywordType_CONJUNCTION,

			"E ": messages.StepKeywordType_CONJUNCTION,

			"Pero ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"et": &Dialect{
		Language: "et",
		Name:     "Estonian",
		Native:   "eesti keel",
		Keywords: map[string][]string{
			feature: {
				"Omadus",
			},
			rule: {
				"Reegel",
			},
			background: {
				"Taust",
			},
			scenario: {
				"Juhtum",
				"Stsenaarium",
			},
			scenarioOutline: {
				"Raamjuhtum",
				"Raamstsenaarium",
			},
			examples: {
				"Juhtumid",
			},
			given: {
				"* ",
				"Eeldades ",
			},
			when: {
				"* ",
				"Kui ",
			},
			then: {
				"* ",
				"Siis ",
			},
			and: {
				"* ",
				"Ja ",
			},
			but: {
				"* ",
				"Kuid ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Eeldades ": messages.StepKeywordType_CONTEXT,

			"Kui ": messages.StepKeywordType_ACTION,

			"Siis ": messages.StepKeywordType_OUTCOME,

			"Ja ": messages.StepKeywordType_CONJUNCTION,

			"Kuid ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"fa": &Dialect{
		Language: "fa",
		Name:     "Persian",
		Native:   "فارسی",
		Keywords: map[string][]string{
			feature: {
				"ویژگی",
				"قابلیت",
			},
			rule: {
				"قانون",
			},
			background: {
				"زمینه",
				"پیش زمینه",
				"مقدمات",
			},
			scenario: {
				"مثال",
				"سناریو",
			},
			scenarioOutline: {
				"الگوی سناریو",
			},
			examples: {
				"نمونه ها",
			},
			given: {
				"* ",
				"با فرض ",
				"فرض کنید ",
				"با در نظر گرفتن ",
			},
			when: {
				"* ",
				"هنگامی ",
				"وقتی ",
			},
			then: {
				"* ",
				"آنگاه ",
				"سپس ",
				"انتظار می رود ",
			},
			and: {
				"* ",
				"و ",
			},
			but: {
				"* ",
				"اما ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"با فرض ": messages.StepKeywordType_CONTEXT,

			"فرض کنید ": messages.StepKeywordType_CONTEXT,

			"با در نظر گرفتن ": messages.StepKeywordType_CONTEXT,

			"هنگامی ": messages.StepKeywordType_ACTION,

			"وقتی ": messages.StepKeywordType_ACTION,

			"آنگاه ": messages.StepKeywordType_OUTCOME,

			"سپس ": messages.StepKeywordType_OUTCOME,

			"انتظار می رود ": messages.StepKeywordType_OUTCOME,

			"و ": messages.StepKeywordType_CONJUNCTION,

			"اما ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"fi": &Dialect{
		Language: "fi",
		Name:     "Finnish",
		Native:   "suomi",
		Keywords: map[string][]string{
			feature: {
				"Ominaisuus",
			},
			rule: {
				"Rule",
			},
			background: {
				"Tausta",
			},
			scenario: {
				"Tapaus",
			},
			scenarioOutline: {
				"Tapausaihio",
			},
			examples: {
				"Tapaukset",
			},
			given: {
				"* ",
				"Oletetaan ",
			},
			when: {
				"* ",
				"Kun ",
			},
			then: {
				"* ",
				"Niin ",
			},
			and: {
				"* ",
				"Ja ",
			},
			but: {
				"* ",
				"Mutta ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Oletetaan ": messages.StepKeywordType_CONTEXT,

			"Kun ": messages.StepKeywordType_ACTION,

			"Niin ": messages.StepKeywordType_OUTCOME,

			"Ja ": messages.StepKeywordType_CONJUNCTION,

			"Mutta ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"fr": &Dialect{
		Language: "fr",
		Name:     "French",
		Native:   "français",
		Keywords: map[string][]string{
			feature: {
				"Fonctionnalité",
			},
			rule: {
				"Règle",
			},
			background: {
				"Contexte",
			},
			scenario: {
				"Exemple",
				"Scénario",
			},
			scenarioOutline: {
				"Plan du scénario",
				"Plan du Scénario",
			},
			examples: {
				"Exemples",
			},
			given: {
				"* ",
				"Soit ",
				"Sachant que ",
				"Sachant qu'",
				"Sachant ",
				"Etant donné que ",
				"Etant donné qu'",
				"Etant donné ",
				"Etant donnée ",
				"Etant donnés ",
				"Etant données ",
				"Étant donné que ",
				"Étant donné qu'",
				"Étant donné ",
				"Étant donnée ",
				"Étant donnés ",
				"Étant données ",
			},
			when: {
				"* ",
				"Quand ",
				"Lorsque ",
				"Lorsqu'",
			},
			then: {
				"* ",
				"Alors ",
				"Donc ",
			},
			and: {
				"* ",
				"Et que ",
				"Et qu'",
				"Et ",
			},
			but: {
				"* ",
				"Mais que ",
				"Mais qu'",
				"Mais ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Soit ": messages.StepKeywordType_CONTEXT,

			"Sachant que ": messages.StepKeywordType_CONTEXT,

			"Sachant qu'": messages.StepKeywordType_CONTEXT,

			"Sachant ": messages.StepKeywordType_CONTEXT,

			"Etant donné que ": messages.StepKeywordType_CONTEXT,

			"Etant donné qu'": messages.StepKeywordType_CONTEXT,

			"Etant donné ": messages.StepKeywordType_CONTEXT,

			"Etant donnée ": messages.StepKeywordType_CONTEXT,

			"Etant donnés ": messages.StepKeywordType_CONTEXT,

			"Etant données ": messages.StepKeywordType_CONTEXT,

			"Étant donné que ": messages.StepKeywordType_CONTEXT,

			"Étant donné qu'": messages.StepKeywordType_CONTEXT,

			"Étant donné ": messages.StepKeywordType_CONTEXT,

			"Étant donnée ": messages.StepKeywordType_CONTEXT,

			"Étant donnés ": messages.StepKeywordType_CONTEXT,

			"Étant données ": messages.StepKeywordType_CONTEXT,

			"Quand ": messages.StepKeywordType_ACTION,

			"Lorsque ": messages.StepKeywordType_ACTION,

			"Lorsqu'": messages.StepKeywordType_ACTION,

			"Alors ": messages.StepKeywordType_OUTCOME,

			"Donc ": messages.StepKeywordType_OUTCOME,

			"Et que ": messages.StepKeywordType_CONJUNCTION,

			"Et qu'": messages.StepKeywordType_CONJUNCTION,

			"Et ": messages.StepKeywordType_CONJUNCTION,

			"Mais que ": messages.StepKeywordType_CONJUNCTION,

			"Mais qu'": messages.StepKeywordType_CONJUNCTION,

			"Mais ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ga": &Dialect{
		Language: "ga",
		Name:     "Irish",
		Native:   "Gaeilge",
		Keywords: map[string][]string{
			feature: {
				"Gné",
			},
			rule: {
				"Riail",
			},
			background: {
				"Cúlra",
			},
			scenario: {
				"Sampla",
				"Cás",
			},
			scenarioOutline: {
				"Cás Achomair",
			},
			examples: {
				"Samplaí",
			},
			given: {
				"* ",
				"Cuir i gcás go ",
				"Cuir i gcás nach ",
				"Cuir i gcás gur ",
				"Cuir i gcás nár ",
			},
			when: {
				"* ",
				"Nuair a ",
				"Nuair nach ",
				"Nuair ba ",
				"Nuair nár ",
			},
			then: {
				"* ",
				"Ansin ",
			},
			and: {
				"* ",
				"Agus ",
			},
			but: {
				"* ",
				"Ach ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Cuir i gcás go ": messages.StepKeywordType_CONTEXT,

			"Cuir i gcás nach ": messages.StepKeywordType_CONTEXT,

			"Cuir i gcás gur ": messages.StepKeywordType_CONTEXT,

			"Cuir i gcás nár ": messages.StepKeywordType_CONTEXT,

			"Nuair a ": messages.StepKeywordType_ACTION,

			"Nuair nach ": messages.StepKeywordType_ACTION,

			"Nuair ba ": messages.StepKeywordType_ACTION,

			"Nuair nár ": messages.StepKeywordType_ACTION,

			"Ansin ": messages.StepKeywordType_OUTCOME,

			"Agus ": messages.StepKeywordType_CONJUNCTION,

			"Ach ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"gj": &Dialect{
		Language: "gj",
		Name:     "Gujarati",
		Native:   "ગુજરાતી",
		Keywords: map[string][]string{
			feature: {
				"લક્ષણ",
				"વ્યાપાર જરૂર",
				"ક્ષમતા",
			},
			rule: {
				"નિયમ",
			},
			background: {
				"બેકગ્રાઉન્ડ",
			},
			scenario: {
				"ઉદાહરણ",
				"સ્થિતિ",
			},
			scenarioOutline: {
				"પરિદ્દશ્ય રૂપરેખા",
				"પરિદ્દશ્ય ઢાંચો",
			},
			examples: {
				"ઉદાહરણો",
			},
			given: {
				"* ",
				"આપેલ છે ",
			},
			when: {
				"* ",
				"ક્યારે ",
			},
			then: {
				"* ",
				"પછી ",
			},
			and: {
				"* ",
				"અને ",
			},
			but: {
				"* ",
				"પણ ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"આપેલ છે ": messages.StepKeywordType_CONTEXT,

			"ક્યારે ": messages.StepKeywordType_ACTION,

			"પછી ": messages.StepKeywordType_OUTCOME,

			"અને ": messages.StepKeywordType_CONJUNCTION,

			"પણ ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"gl": &Dialect{
		Language: "gl",
		Name:     "Galician",
		Native:   "galego",
		Keywords: map[string][]string{
			feature: {
				"Característica",
			},
			rule: {
				"Rule",
			},
			background: {
				"Contexto",
			},
			scenario: {
				"Exemplo",
				"Escenario",
			},
			scenarioOutline: {
				"Esbozo do escenario",
			},
			examples: {
				"Exemplos",
			},
			given: {
				"* ",
				"Dado ",
				"Dada ",
				"Dados ",
				"Dadas ",
			},
			when: {
				"* ",
				"Cando ",
			},
			then: {
				"* ",
				"Entón ",
				"Logo ",
			},
			and: {
				"* ",
				"E ",
			},
			but: {
				"* ",
				"Mais ",
				"Pero ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dado ": messages.StepKeywordType_CONTEXT,

			"Dada ": messages.StepKeywordType_CONTEXT,

			"Dados ": messages.StepKeywordType_CONTEXT,

			"Dadas ": messages.StepKeywordType_CONTEXT,

			"Cando ": messages.StepKeywordType_ACTION,

			"Entón ": messages.StepKeywordType_OUTCOME,

			"Logo ": messages.StepKeywordType_OUTCOME,

			"E ": messages.StepKeywordType_CONJUNCTION,

			"Mais ": messages.StepKeywordType_CONJUNCTION,

			"Pero ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"he": &Dialect{
		Language: "he",
		Name:     "Hebrew",
		Native:   "עברית",
		Keywords: map[string][]string{
			feature: {
				"תכונה",
			},
			rule: {
				"כלל",
			},
			background: {
				"רקע",
			},
			scenario: {
				"דוגמא",
				"תרחיש",
			},
			scenarioOutline: {
				"תבנית תרחיש",
			},
			examples: {
				"דוגמאות",
			},
			given: {
				"* ",
				"בהינתן ",
			},
			when: {
				"* ",
				"כאשר ",
			},
			then: {
				"* ",
				"אז ",
				"אזי ",
			},
			and: {
				"* ",
				"וגם ",
			},
			but: {
				"* ",
				"אבל ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"בהינתן ": messages.StepKeywordType_CONTEXT,

			"כאשר ": messages.StepKeywordType_ACTION,

			"אז ": messages.StepKeywordType_OUTCOME,

			"אזי ": messages.StepKeywordType_OUTCOME,

			"וגם ": messages.StepKeywordType_CONJUNCTION,

			"אבל ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"hi": &Dialect{
		Language: "hi",
		Name:     "Hindi",
		Native:   "हिंदी",
		Keywords: map[string][]string{
			feature: {
				"रूप लेख",
			},
			rule: {
				"नियम",
			},
			background: {
				"पृष्ठभूमि",
			},
			scenario: {
				"परिदृश्य",
			},
			scenarioOutline: {
				"परिदृश्य रूपरेखा",
			},
			examples: {
				"उदाहरण",
			},
			given: {
				"* ",
				"अगर ",
				"यदि ",
				"चूंकि ",
			},
			when: {
				"* ",
				"जब ",
				"कदा ",
			},
			then: {
				"* ",
				"तब ",
				"तदा ",
			},
			and: {
				"* ",
				"और ",
				"तथा ",
			},
			but: {
				"* ",
				"पर ",
				"परन्तु ",
				"किन्तु ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"अगर ": messages.StepKeywordType_CONTEXT,

			"यदि ": messages.StepKeywordType_CONTEXT,

			"चूंकि ": messages.StepKeywordType_CONTEXT,

			"जब ": messages.StepKeywordType_ACTION,

			"कदा ": messages.StepKeywordType_ACTION,

			"तब ": messages.StepKeywordType_OUTCOME,

			"तदा ": messages.StepKeywordType_OUTCOME,

			"और ": messages.StepKeywordType_CONJUNCTION,

			"तथा ": messages.StepKeywordType_CONJUNCTION,

			"पर ": messages.StepKeywordType_CONJUNCTION,

			"परन्तु ": messages.StepKeywordType_CONJUNCTION,

			"किन्तु ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"hr": &Dialect{
		Language: "hr",
		Name:     "Croatian",
		Native:   "hrvatski",
		Keywords: map[string][]string{
			feature: {
				"Osobina",
				"Mogućnost",
				"Mogucnost",
			},
			rule: {
				"Rule",
			},
			background: {
				"Pozadina",
			},
			scenario: {
				"Primjer",
				"Scenarij",
			},
			scenarioOutline: {
				"Skica",
				"Koncept",
			},
			examples: {
				"Primjeri",
				"Scenariji",
			},
			given: {
				"* ",
				"Zadan ",
				"Zadani ",
				"Zadano ",
				"Ukoliko ",
			},
			when: {
				"* ",
				"Kada ",
				"Kad ",
			},
			then: {
				"* ",
				"Onda ",
			},
			and: {
				"* ",
				"I ",
			},
			but: {
				"* ",
				"Ali ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Zadan ": messages.StepKeywordType_CONTEXT,

			"Zadani ": messages.StepKeywordType_CONTEXT,

			"Zadano ": messages.StepKeywordType_CONTEXT,

			"Ukoliko ": messages.StepKeywordType_CONTEXT,

			"Kada ": messages.StepKeywordType_ACTION,

			"Kad ": messages.StepKeywordType_ACTION,

			"Onda ": messages.StepKeywordType_OUTCOME,

			"I ": messages.StepKeywordType_CONJUNCTION,

			"Ali ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ht": &Dialect{
		Language: "ht",
		Name:     "Creole",
		Native:   "kreyòl",
		Keywords: map[string][]string{
			feature: {
				"Karakteristik",
				"Mak",
				"Fonksyonalite",
			},
			rule: {
				"Rule",
			},
			background: {
				"Kontèks",
				"Istorik",
			},
			scenario: {
				"Senaryo",
			},
			scenarioOutline: {
				"Plan senaryo",
				"Plan Senaryo",
				"Senaryo deskripsyon",
				"Senaryo Deskripsyon",
				"Dyagram senaryo",
				"Dyagram Senaryo",
			},
			examples: {
				"Egzanp",
			},
			given: {
				"* ",
				"Sipoze ",
				"Sipoze ke ",
				"Sipoze Ke ",
			},
			when: {
				"* ",
				"Lè ",
				"Le ",
			},
			then: {
				"* ",
				"Lè sa a ",
				"Le sa a ",
			},
			and: {
				"* ",
				"Ak ",
				"Epi ",
				"E ",
			},
			but: {
				"* ",
				"Men ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Sipoze ": messages.StepKeywordType_CONTEXT,

			"Sipoze ke ": messages.StepKeywordType_CONTEXT,

			"Sipoze Ke ": messages.StepKeywordType_CONTEXT,

			"Lè ": messages.StepKeywordType_ACTION,

			"Le ": messages.StepKeywordType_ACTION,

			"Lè sa a ": messages.StepKeywordType_OUTCOME,

			"Le sa a ": messages.StepKeywordType_OUTCOME,

			"Ak ": messages.StepKeywordType_CONJUNCTION,

			"Epi ": messages.StepKeywordType_CONJUNCTION,

			"E ": messages.StepKeywordType_CONJUNCTION,

			"Men ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"hu": &Dialect{
		Language: "hu",
		Name:     "Hungarian",
		Native:   "magyar",
		Keywords: map[string][]string{
			feature: {
				"Jellemző",
			},
			rule: {
				"Szabály",
			},
			background: {
				"Háttér",
			},
			scenario: {
				"Példa",
				"Forgatókönyv",
			},
			scenarioOutline: {
				"Forgatókönyv vázlat",
			},
			examples: {
				"Példák",
			},
			given: {
				"* ",
				"Amennyiben ",
				"Adott ",
			},
			when: {
				"* ",
				"Majd ",
				"Ha ",
				"Amikor ",
			},
			then: {
				"* ",
				"Akkor ",
			},
			and: {
				"* ",
				"És ",
			},
			but: {
				"* ",
				"De ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Amennyiben ": messages.StepKeywordType_CONTEXT,

			"Adott ": messages.StepKeywordType_CONTEXT,

			"Majd ": messages.StepKeywordType_ACTION,

			"Ha ": messages.StepKeywordType_ACTION,

			"Amikor ": messages.StepKeywordType_ACTION,

			"Akkor ": messages.StepKeywordType_OUTCOME,

			"És ": messages.StepKeywordType_CONJUNCTION,

			"De ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"id": &Dialect{
		Language: "id",
		Name:     "Indonesian",
		Native:   "Bahasa Indonesia",
		Keywords: map[string][]string{
			feature: {
				"Fitur",
			},
			rule: {
				"Rule",
				"Aturan",
			},
			background: {
				"Dasar",
				"Latar Belakang",
			},
			scenario: {
				"Skenario",
			},
			scenarioOutline: {
				"Skenario konsep",
				"Garis-Besar Skenario",
			},
			examples: {
				"Contoh",
				"Misal",
			},
			given: {
				"* ",
				"Dengan ",
				"Diketahui ",
				"Diasumsikan ",
				"Bila ",
				"Jika ",
			},
			when: {
				"* ",
				"Ketika ",
			},
			then: {
				"* ",
				"Maka ",
				"Kemudian ",
			},
			and: {
				"* ",
				"Dan ",
			},
			but: {
				"* ",
				"Tapi ",
				"Tetapi ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dengan ": messages.StepKeywordType_CONTEXT,

			"Diketahui ": messages.StepKeywordType_CONTEXT,

			"Diasumsikan ": messages.StepKeywordType_CONTEXT,

			"Bila ": messages.StepKeywordType_CONTEXT,

			"Jika ": messages.StepKeywordType_CONTEXT,

			"Ketika ": messages.StepKeywordType_ACTION,

			"Maka ": messages.StepKeywordType_OUTCOME,

			"Kemudian ": messages.StepKeywordType_OUTCOME,

			"Dan ": messages.StepKeywordType_CONJUNCTION,

			"Tapi ": messages.StepKeywordType_CONJUNCTION,

			"Tetapi ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"is": &Dialect{
		Language: "is",
		Name:     "Icelandic",
		Native:   "Íslenska",
		Keywords: map[string][]string{
			feature: {
				"Eiginleiki",
			},
			rule: {
				"Rule",
			},
			background: {
				"Bakgrunnur",
			},
			scenario: {
				"Atburðarás",
			},
			scenarioOutline: {
				"Lýsing Atburðarásar",
				"Lýsing Dæma",
			},
			examples: {
				"Dæmi",
				"Atburðarásir",
			},
			given: {
				"* ",
				"Ef ",
			},
			when: {
				"* ",
				"Þegar ",
			},
			then: {
				"* ",
				"Þá ",
			},
			and: {
				"* ",
				"Og ",
			},
			but: {
				"* ",
				"En ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Ef ": messages.StepKeywordType_CONTEXT,

			"Þegar ": messages.StepKeywordType_ACTION,

			"Þá ": messages.StepKeywordType_OUTCOME,

			"Og ": messages.StepKeywordType_CONJUNCTION,

			"En ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"it": &Dialect{
		Language: "it",
		Name:     "Italian",
		Native:   "italiano",
		Keywords: map[string][]string{
			feature: {
				"Funzionalità",
				"Esigenza di Business",
				"Abilità",
			},
			rule: {
				"Regola",
			},
			background: {
				"Contesto",
			},
			scenario: {
				"Esempio",
				"Scenario",
			},
			scenarioOutline: {
				"Schema dello scenario",
			},
			examples: {
				"Esempi",
			},
			given: {
				"* ",
				"Dato ",
				"Data ",
				"Dati ",
				"Date ",
			},
			when: {
				"* ",
				"Quando ",
			},
			then: {
				"* ",
				"Allora ",
			},
			and: {
				"* ",
				"E ",
				"Ed ",
			},
			but: {
				"* ",
				"Ma ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dato ": messages.StepKeywordType_CONTEXT,

			"Data ": messages.StepKeywordType_CONTEXT,

			"Dati ": messages.StepKeywordType_CONTEXT,

			"Date ": messages.StepKeywordType_CONTEXT,

			"Quando ": messages.StepKeywordType_ACTION,

			"Allora ": messages.StepKeywordType_OUTCOME,

			"E ": messages.StepKeywordType_CONJUNCTION,

			"Ed ": messages.StepKeywordType_CONJUNCTION,

			"Ma ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ja": &Dialect{
		Language: "ja",
		Name:     "Japanese",
		Native:   "日本語",
		Keywords: map[string][]string{
			feature: {
				"フィーチャ",
				"機能",
			},
			rule: {
				"ルール",
			},
			background: {
				"背景",
			},
			scenario: {
				"シナリオ",
			},
			scenarioOutline: {
				"シナリオアウトライン",
				"シナリオテンプレート",
				"テンプレ",
				"シナリオテンプレ",
			},
			examples: {
				"例",
				"サンプル",
			},
			given: {
				"* ",
				"前提",
			},
			when: {
				"* ",
				"もし",
			},
			then: {
				"* ",
				"ならば",
			},
			and: {
				"* ",
				"且つ",
				"かつ",
			},
			but: {
				"* ",
				"然し",
				"しかし",
				"但し",
				"ただし",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"前提": messages.StepKeywordType_CONTEXT,

			"もし": messages.StepKeywordType_ACTION,

			"ならば": messages.StepKeywordType_OUTCOME,

			"且つ": messages.StepKeywordType_CONJUNCTION,

			"かつ": messages.StepKeywordType_CONJUNCTION,

			"然し": messages.StepKeywordType_CONJUNCTION,

			"しかし": messages.StepKeywordType_CONJUNCTION,

			"但し": messages.StepKeywordType_CONJUNCTION,

			"ただし": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"jv": &Dialect{
		Language: "jv",
		Name:     "Javanese",
		Native:   "Basa Jawa",
		Keywords: map[string][]string{
			feature: {
				"Fitur",
			},
			rule: {
				"Rule",
			},
			background: {
				"Dasar",
			},
			scenario: {
				"Skenario",
			},
			scenarioOutline: {
				"Konsep skenario",
			},
			examples: {
				"Conto",
				"Contone",
			},
			given: {
				"* ",
				"Nalika ",
				"Nalikaning ",
			},
			when: {
				"* ",
				"Manawa ",
				"Menawa ",
			},
			then: {
				"* ",
				"Njuk ",
				"Banjur ",
			},
			and: {
				"* ",
				"Lan ",
			},
			but: {
				"* ",
				"Tapi ",
				"Nanging ",
				"Ananging ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Nalika ": messages.StepKeywordType_CONTEXT,

			"Nalikaning ": messages.StepKeywordType_CONTEXT,

			"Manawa ": messages.StepKeywordType_ACTION,

			"Menawa ": messages.StepKeywordType_ACTION,

			"Njuk ": messages.StepKeywordType_OUTCOME,

			"Banjur ": messages.StepKeywordType_OUTCOME,

			"Lan ": messages.StepKeywordType_CONJUNCTION,

			"Tapi ": messages.StepKeywordType_CONJUNCTION,

			"Nanging ": messages.StepKeywordType_CONJUNCTION,

			"Ananging ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ka": &Dialect{
		Language: "ka",
		Name:     "Georgian",
		Native:   "ქართული",
		Keywords: map[string][]string{
			feature: {
				"თვისება",
				"მოთხოვნა",
			},
			rule: {
				"წესი",
			},
			background: {
				"კონტექსტი",
			},
			scenario: {
				"მაგალითად",
				"მაგალითი",
				"მაგ",
				"სცენარი",
			},
			scenarioOutline: {
				"სცენარის ნიმუში",
				"სცენარის შაბლონი",
				"ნიმუში",
				"შაბლონი",
			},
			examples: {
				"მაგალითები",
			},
			given: {
				"* ",
				"მოცემული ",
				"მოცემულია ",
				"ვთქვათ ",
			},
			when: {
				"* ",
				"როდესაც ",
				"როცა ",
				"როგორც კი ",
				"თუ ",
			},
			then: {
				"* ",
				"მაშინ ",
			},
			and: {
				"* ",
				"და ",
				"ასევე ",
			},
			but: {
				"* ",
				"მაგრამ ",
				"თუმცა ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"მოცემული ": messages.StepKeywordType_CONTEXT,

			"მოცემულია ": messages.StepKeywordType_CONTEXT,

			"ვთქვათ ": messages.StepKeywordType_CONTEXT,

			"როდესაც ": messages.StepKeywordType_ACTION,

			"როცა ": messages.StepKeywordType_ACTION,

			"როგორც კი ": messages.StepKeywordType_ACTION,

			"თუ ": messages.StepKeywordType_ACTION,

			"მაშინ ": messages.StepKeywordType_OUTCOME,

			"და ": messages.StepKeywordType_CONJUNCTION,

			"ასევე ": messages.StepKeywordType_CONJUNCTION,

			"მაგრამ ": messages.StepKeywordType_CONJUNCTION,

			"თუმცა ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"kn": &Dialect{
		Language: "kn",
		Name:     "Kannada",
		Native:   "ಕನ್ನಡ",
		Keywords: map[string][]string{
			feature: {
				"ಹೆಚ್ಚಳ",
			},
			rule: {
				"Rule",
			},
			background: {
				"ಹಿನ್ನೆಲೆ",
			},
			scenario: {
				"ಉದಾಹರಣೆ",
				"ಕಥಾಸಾರಾಂಶ",
			},
			scenarioOutline: {
				"ವಿವರಣೆ",
			},
			examples: {
				"ಉದಾಹರಣೆಗಳು",
			},
			given: {
				"* ",
				"ನೀಡಿದ ",
			},
			when: {
				"* ",
				"ಸ್ಥಿತಿಯನ್ನು ",
			},
			then: {
				"* ",
				"ನಂತರ ",
			},
			and: {
				"* ",
				"ಮತ್ತು ",
			},
			but: {
				"* ",
				"ಆದರೆ ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"ನೀಡಿದ ": messages.StepKeywordType_CONTEXT,

			"ಸ್ಥಿತಿಯನ್ನು ": messages.StepKeywordType_ACTION,

			"ನಂತರ ": messages.StepKeywordType_OUTCOME,

			"ಮತ್ತು ": messages.StepKeywordType_CONJUNCTION,

			"ಆದರೆ ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ko": &Dialect{
		Language: "ko",
		Name:     "Korean",
		Native:   "한국어",
		Keywords: map[string][]string{
			feature: {
				"기능",
			},
			rule: {
				"규칙",
			},
			background: {
				"배경",
			},
			scenario: {
				"시나리오",
			},
			scenarioOutline: {
				"시나리오 개요",
			},
			examples: {
				"예",
			},
			given: {
				"* ",
				"조건 ",
				"먼저 ",
			},
			when: {
				"* ",
				"만일 ",
				"만약 ",
			},
			then: {
				"* ",
				"그러면 ",
			},
			and: {
				"* ",
				"그리고 ",
			},
			but: {
				"* ",
				"하지만 ",
				"단 ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"조건 ": messages.StepKeywordType_CONTEXT,

			"먼저 ": messages.StepKeywordType_CONTEXT,

			"만일 ": messages.StepKeywordType_ACTION,

			"만약 ": messages.StepKeywordType_ACTION,

			"그러면 ": messages.StepKeywordType_OUTCOME,

			"그리고 ": messages.StepKeywordType_CONJUNCTION,

			"하지만 ": messages.StepKeywordType_CONJUNCTION,

			"단 ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"lt": &Dialect{
		Language: "lt",
		Name:     "Lithuanian",
		Native:   "lietuvių kalba",
		Keywords: map[string][]string{
			feature: {
				"Savybė",
			},
			rule: {
				"Rule",
			},
			background: {
				"Kontekstas",
			},
			scenario: {
				"Pavyzdys",
				"Scenarijus",
			},
			scenarioOutline: {
				"Scenarijaus šablonas",
			},
			examples: {
				"Pavyzdžiai",
				"Scenarijai",
				"Variantai",
			},
			given: {
				"* ",
				"Duota ",
			},
			when: {
				"* ",
				"Kai ",
			},
			then: {
				"* ",
				"Tada ",
			},
			and: {
				"* ",
				"Ir ",
			},
			but: {
				"* ",
				"Bet ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Duota ": messages.StepKeywordType_CONTEXT,

			"Kai ": messages.StepKeywordType_ACTION,

			"Tada ": messages.StepKeywordType_OUTCOME,

			"Ir ": messages.StepKeywordType_CONJUNCTION,

			"Bet ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"lu": &Dialect{
		Language: "lu",
		Name:     "Luxemburgish",
		Native:   "Lëtzebuergesch",
		Keywords: map[string][]string{
			feature: {
				"Funktionalitéit",
			},
			rule: {
				"Rule",
			},
			background: {
				"Hannergrond",
			},
			scenario: {
				"Beispill",
				"Szenario",
			},
			scenarioOutline: {
				"Plang vum Szenario",
			},
			examples: {
				"Beispiller",
			},
			given: {
				"* ",
				"ugeholl ",
			},
			when: {
				"* ",
				"wann ",
			},
			then: {
				"* ",
				"dann ",
			},
			and: {
				"* ",
				"an ",
				"a ",
			},
			but: {
				"* ",
				"awer ",
				"mä ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"ugeholl ": messages.StepKeywordType_CONTEXT,

			"wann ": messages.StepKeywordType_ACTION,

			"dann ": messages.StepKeywordType_OUTCOME,

			"an ": messages.StepKeywordType_CONJUNCTION,

			"a ": messages.StepKeywordType_CONJUNCTION,

			"awer ": messages.StepKeywordType_CONJUNCTION,

			"mä ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"lv": &Dialect{
		Language: "lv",
		Name:     "Latvian",
		Native:   "latviešu",
		Keywords: map[string][]string{
			feature: {
				"Funkcionalitāte",
				"Fīča",
			},
			rule: {
				"Rule",
			},
			background: {
				"Konteksts",
				"Situācija",
			},
			scenario: {
				"Piemērs",
				"Scenārijs",
			},
			scenarioOutline: {
				"Scenārijs pēc parauga",
			},
			examples: {
				"Piemēri",
				"Paraugs",
			},
			given: {
				"* ",
				"Kad ",
			},
			when: {
				"* ",
				"Ja ",
			},
			then: {
				"* ",
				"Tad ",
			},
			and: {
				"* ",
				"Un ",
			},
			but: {
				"* ",
				"Bet ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Kad ": messages.StepKeywordType_CONTEXT,

			"Ja ": messages.StepKeywordType_ACTION,

			"Tad ": messages.StepKeywordType_OUTCOME,

			"Un ": messages.StepKeywordType_CONJUNCTION,

			"Bet ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"mk-Cyrl": &Dialect{
		Language: "mk-Cyrl",
		Name:     "Macedonian",
		Native:   "Македонски",
		Keywords: map[string][]string{
			feature: {
				"Функционалност",
				"Бизнис потреба",
				"Можност",
			},
			rule: {
				"Rule",
			},
			background: {
				"Контекст",
				"Содржина",
			},
			scenario: {
				"Пример",
				"Сценарио",
				"На пример",
			},
			scenarioOutline: {
				"Преглед на сценарија",
				"Скица",
				"Концепт",
			},
			examples: {
				"Примери",
				"Сценарија",
			},
			given: {
				"* ",
				"Дадено ",
				"Дадена ",
			},
			when: {
				"* ",
				"Кога ",
			},
			then: {
				"* ",
				"Тогаш ",
			},
			and: {
				"* ",
				"И ",
			},
			but: {
				"* ",
				"Но ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Дадено ": messages.StepKeywordType_CONTEXT,

			"Дадена ": messages.StepKeywordType_CONTEXT,

			"Кога ": messages.StepKeywordType_ACTION,

			"Тогаш ": messages.StepKeywordType_OUTCOME,

			"И ": messages.StepKeywordType_CONJUNCTION,

			"Но ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"mk-Latn": &Dialect{
		Language: "mk-Latn",
		Name:     "Macedonian (Latin)",
		Native:   "Makedonski (Latinica)",
		Keywords: map[string][]string{
			feature: {
				"Funkcionalnost",
				"Biznis potreba",
				"Mozhnost",
			},
			rule: {
				"Rule",
			},
			background: {
				"Kontekst",
				"Sodrzhina",
			},
			scenario: {
				"Scenario",
				"Na primer",
			},
			scenarioOutline: {
				"Pregled na scenarija",
				"Skica",
				"Koncept",
			},
			examples: {
				"Primeri",
				"Scenaria",
			},
			given: {
				"* ",
				"Dadeno ",
				"Dadena ",
			},
			when: {
				"* ",
				"Koga ",
			},
			then: {
				"* ",
				"Togash ",
			},
			and: {
				"* ",
				"I ",
			},
			but: {
				"* ",
				"No ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dadeno ": messages.StepKeywordType_CONTEXT,

			"Dadena ": messages.StepKeywordType_CONTEXT,

			"Koga ": messages.StepKeywordType_ACTION,

			"Togash ": messages.StepKeywordType_OUTCOME,

			"I ": messages.StepKeywordType_CONJUNCTION,

			"No ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"mn": &Dialect{
		Language: "mn",
		Name:     "Mongolian",
		Native:   "монгол",
		Keywords: map[string][]string{
			feature: {
				"Функц",
				"Функционал",
			},
			rule: {
				"Rule",
			},
			background: {
				"Агуулга",
			},
			scenario: {
				"Сценар",
			},
			scenarioOutline: {
				"Сценарын төлөвлөгөө",
			},
			examples: {
				"Тухайлбал",
			},
			given: {
				"* ",
				"Өгөгдсөн нь ",
				"Анх ",
			},
			when: {
				"* ",
				"Хэрэв ",
			},
			then: {
				"* ",
				"Тэгэхэд ",
				"Үүний дараа ",
			},
			and: {
				"* ",
				"Мөн ",
				"Тэгээд ",
			},
			but: {
				"* ",
				"Гэхдээ ",
				"Харин ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Өгөгдсөн нь ": messages.StepKeywordType_CONTEXT,

			"Анх ": messages.StepKeywordType_CONTEXT,

			"Хэрэв ": messages.StepKeywordType_ACTION,

			"Тэгэхэд ": messages.StepKeywordType_OUTCOME,

			"Үүний дараа ": messages.StepKeywordType_OUTCOME,

			"Мөн ": messages.StepKeywordType_CONJUNCTION,

			"Тэгээд ": messages.StepKeywordType_CONJUNCTION,

			"Гэхдээ ": messages.StepKeywordType_CONJUNCTION,

			"Харин ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ne": &Dialect{
		Language: "ne",
		Name:     "Nepali",
		Native:   "नेपाली",
		Keywords: map[string][]string{
			feature: {
				"सुविधा",
				"विशेषता",
			},
			rule: {
				"नियम",
			},
			background: {
				"पृष्ठभूमी",
			},
			scenario: {
				"परिदृश्य",
			},
			scenarioOutline: {
				"परिदृश्य रूपरेखा",
			},
			examples: {
				"उदाहरण",
				"उदाहरणहरु",
			},
			given: {
				"* ",
				"दिइएको ",
				"दिएको ",
				"यदि ",
			},
			when: {
				"* ",
				"जब ",
			},
			then: {
				"* ",
				"त्यसपछि ",
				"अनी ",
			},
			and: {
				"* ",
				"र ",
				"अनि ",
			},
			but: {
				"* ",
				"तर ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"दिइएको ": messages.StepKeywordType_CONTEXT,

			"दिएको ": messages.StepKeywordType_CONTEXT,

			"यदि ": messages.StepKeywordType_CONTEXT,

			"जब ": messages.StepKeywordType_ACTION,

			"त्यसपछि ": messages.StepKeywordType_OUTCOME,

			"अनी ": messages.StepKeywordType_OUTCOME,

			"र ": messages.StepKeywordType_CONJUNCTION,

			"अनि ": messages.StepKeywordType_CONJUNCTION,

			"तर ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"nl": &Dialect{
		Language: "nl",
		Name:     "Dutch",
		Native:   "Nederlands",
		Keywords: map[string][]string{
			feature: {
				"Functionaliteit",
			},
			rule: {
				"Regel",
			},
			background: {
				"Achtergrond",
			},
			scenario: {
				"Voorbeeld",
				"Scenario",
			},
			scenarioOutline: {
				"Abstract Scenario",
			},
			examples: {
				"Voorbeelden",
			},
			given: {
				"* ",
				"Gegeven ",
				"Stel ",
			},
			when: {
				"* ",
				"Als ",
				"Wanneer ",
			},
			then: {
				"* ",
				"Dan ",
			},
			and: {
				"* ",
				"En ",
			},
			but: {
				"* ",
				"Maar ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Gegeven ": messages.StepKeywordType_CONTEXT,

			"Stel ": messages.StepKeywordType_CONTEXT,

			"Als ": messages.StepKeywordType_ACTION,

			"Wanneer ": messages.StepKeywordType_ACTION,

			"Dan ": messages.StepKeywordType_OUTCOME,

			"En ": messages.StepKeywordType_CONJUNCTION,

			"Maar ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"no": &Dialect{
		Language: "no",
		Name:     "Norwegian",
		Native:   "norsk",
		Keywords: map[string][]string{
			feature: {
				"Egenskap",
			},
			rule: {
				"Regel",
			},
			background: {
				"Bakgrunn",
			},
			scenario: {
				"Eksempel",
				"Scenario",
			},
			scenarioOutline: {
				"Scenariomal",
				"Abstrakt Scenario",
			},
			examples: {
				"Eksempler",
			},
			given: {
				"* ",
				"Gitt ",
			},
			when: {
				"* ",
				"Når ",
			},
			then: {
				"* ",
				"Så ",
			},
			and: {
				"* ",
				"Og ",
			},
			but: {
				"* ",
				"Men ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Gitt ": messages.StepKeywordType_CONTEXT,

			"Når ": messages.StepKeywordType_ACTION,

			"Så ": messages.StepKeywordType_OUTCOME,

			"Og ": messages.StepKeywordType_CONJUNCTION,

			"Men ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"pa": &Dialect{
		Language: "pa",
		Name:     "Panjabi",
		Native:   "ਪੰਜਾਬੀ",
		Keywords: map[string][]string{
			feature: {
				"ਖਾਸੀਅਤ",
				"ਮੁਹਾਂਦਰਾ",
				"ਨਕਸ਼ ਨੁਹਾਰ",
			},
			rule: {
				"Rule",
			},
			background: {
				"ਪਿਛੋਕੜ",
			},
			scenario: {
				"ਉਦਾਹਰਨ",
				"ਪਟਕਥਾ",
			},
			scenarioOutline: {
				"ਪਟਕਥਾ ਢਾਂਚਾ",
				"ਪਟਕਥਾ ਰੂਪ ਰੇਖਾ",
			},
			examples: {
				"ਉਦਾਹਰਨਾਂ",
			},
			given: {
				"* ",
				"ਜੇਕਰ ",
				"ਜਿਵੇਂ ਕਿ ",
			},
			when: {
				"* ",
				"ਜਦੋਂ ",
			},
			then: {
				"* ",
				"ਤਦ ",
			},
			and: {
				"* ",
				"ਅਤੇ ",
			},
			but: {
				"* ",
				"ਪਰ ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"ਜੇਕਰ ": messages.StepKeywordType_CONTEXT,

			"ਜਿਵੇਂ ਕਿ ": messages.StepKeywordType_CONTEXT,

			"ਜਦੋਂ ": messages.StepKeywordType_ACTION,

			"ਤਦ ": messages.StepKeywordType_OUTCOME,

			"ਅਤੇ ": messages.StepKeywordType_CONJUNCTION,

			"ਪਰ ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"pl": &Dialect{
		Language: "pl",
		Name:     "Polish",
		Native:   "polski",
		Keywords: map[string][]string{
			feature: {
				"Właściwość",
				"Funkcja",
				"Aspekt",
				"Potrzeba biznesowa",
			},
			rule: {
				"Zasada",
				"Reguła",
			},
			background: {
				"Założenia",
			},
			scenario: {
				"Przykład",
				"Scenariusz",
			},
			scenarioOutline: {
				"Szablon scenariusza",
			},
			examples: {
				"Przykłady",
			},
			given: {
				"* ",
				"Zakładając ",
				"Mając ",
				"Zakładając, że ",
			},
			when: {
				"* ",
				"Jeżeli ",
				"Jeśli ",
				"Gdy ",
				"Kiedy ",
			},
			then: {
				"* ",
				"Wtedy ",
			},
			and: {
				"* ",
				"Oraz ",
				"I ",
			},
			but: {
				"* ",
				"Ale ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Zakładając ": messages.StepKeywordType_CONTEXT,

			"Mając ": messages.StepKeywordType_CONTEXT,

			"Zakładając, że ": messages.StepKeywordType_CONTEXT,

			"Jeżeli ": messages.StepKeywordType_ACTION,

			"Jeśli ": messages.StepKeywordType_ACTION,

			"Gdy ": messages.StepKeywordType_ACTION,

			"Kiedy ": messages.StepKeywordType_ACTION,

			"Wtedy ": messages.StepKeywordType_OUTCOME,

			"Oraz ": messages.StepKeywordType_CONJUNCTION,

			"I ": messages.StepKeywordType_CONJUNCTION,

			"Ale ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"pt": &Dialect{
		Language: "pt",
		Name:     "Portuguese",
		Native:   "português",
		Keywords: map[string][]string{
			feature: {
				"Funcionalidade",
				"Característica",
				"Caracteristica",
			},
			rule: {
				"Regra",
			},
			background: {
				"Contexto",
				"Cenário de Fundo",
				"Cenario de Fundo",
				"Fundo",
			},
			scenario: {
				"Exemplo",
				"Cenário",
				"Cenario",
			},
			scenarioOutline: {
				"Esquema do Cenário",
				"Esquema do Cenario",
				"Delineação do Cenário",
				"Delineacao do Cenario",
			},
			examples: {
				"Exemplos",
				"Cenários",
				"Cenarios",
			},
			given: {
				"* ",
				"Dado ",
				"Dada ",
				"Dados ",
				"Dadas ",
			},
			when: {
				"* ",
				"Quando ",
			},
			then: {
				"* ",
				"Então ",
				"Entao ",
			},
			and: {
				"* ",
				"E ",
			},
			but: {
				"* ",
				"Mas ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dado ": messages.StepKeywordType_CONTEXT,

			"Dada ": messages.StepKeywordType_CONTEXT,

			"Dados ": messages.StepKeywordType_CONTEXT,

			"Dadas ": messages.StepKeywordType_CONTEXT,

			"Quando ": messages.StepKeywordType_ACTION,

			"Então ": messages.StepKeywordType_OUTCOME,

			"Entao ": messages.StepKeywordType_OUTCOME,

			"E ": messages.StepKeywordType_CONJUNCTION,

			"Mas ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ro": &Dialect{
		Language: "ro",
		Name:     "Romanian",
		Native:   "română",
		Keywords: map[string][]string{
			feature: {
				"Functionalitate",
				"Funcționalitate",
				"Funcţionalitate",
			},
			rule: {
				"Rule",
			},
			background: {
				"Context",
			},
			scenario: {
				"Exemplu",
				"Scenariu",
			},
			scenarioOutline: {
				"Structura scenariu",
				"Structură scenariu",
			},
			examples: {
				"Exemple",
			},
			given: {
				"* ",
				"Date fiind ",
				"Dat fiind ",
				"Dată fiind",
				"Dati fiind ",
				"Dați fiind ",
				"Daţi fiind ",
			},
			when: {
				"* ",
				"Cand ",
				"Când ",
			},
			then: {
				"* ",
				"Atunci ",
			},
			and: {
				"* ",
				"Si ",
				"Și ",
				"Şi ",
			},
			but: {
				"* ",
				"Dar ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Date fiind ": messages.StepKeywordType_CONTEXT,

			"Dat fiind ": messages.StepKeywordType_CONTEXT,

			"Dată fiind": messages.StepKeywordType_CONTEXT,

			"Dati fiind ": messages.StepKeywordType_CONTEXT,

			"Dați fiind ": messages.StepKeywordType_CONTEXT,

			"Daţi fiind ": messages.StepKeywordType_CONTEXT,

			"Cand ": messages.StepKeywordType_ACTION,

			"Când ": messages.StepKeywordType_ACTION,

			"Atunci ": messages.StepKeywordType_OUTCOME,

			"Si ": messages.StepKeywordType_CONJUNCTION,

			"Și ": messages.StepKeywordType_CONJUNCTION,

			"Şi ": messages.StepKeywordType_CONJUNCTION,

			"Dar ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ru": &Dialect{
		Language: "ru",
		Name:     "Russian",
		Native:   "русский",
		Keywords: map[string][]string{
			feature: {
				"Функция",
				"Функциональность",
				"Функционал",
				"Свойство",
				"Фича",
			},
			rule: {
				"Правило",
			},
			background: {
				"Предыстория",
				"Контекст",
			},
			scenario: {
				"Пример",
				"Сценарий",
			},
			scenarioOutline: {
				"Структура сценария",
				"Шаблон сценария",
			},
			examples: {
				"Примеры",
				"Значения",
			},
			given: {
				"* ",
				"Допустим ",
				"Дано ",
				"Пусть ",
			},
			when: {
				"* ",
				"Когда ",
				"Если ",
			},
			then: {
				"* ",
				"То ",
				"Затем ",
				"Тогда ",
			},
			and: {
				"* ",
				"И ",
				"К тому же ",
				"Также ",
			},
			but: {
				"* ",
				"Но ",
				"А ",
				"Иначе ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Допустим ": messages.StepKeywordType_CONTEXT,

			"Дано ": messages.StepKeywordType_CONTEXT,

			"Пусть ": messages.StepKeywordType_CONTEXT,

			"Когда ": messages.StepKeywordType_ACTION,

			"Если ": messages.StepKeywordType_ACTION,

			"То ": messages.StepKeywordType_OUTCOME,

			"Затем ": messages.StepKeywordType_OUTCOME,

			"Тогда ": messages.StepKeywordType_OUTCOME,

			"И ": messages.StepKeywordType_CONJUNCTION,

			"К тому же ": messages.StepKeywordType_CONJUNCTION,

			"Также ": messages.StepKeywordType_CONJUNCTION,

			"Но ": messages.StepKeywordType_CONJUNCTION,

			"А ": messages.StepKeywordType_CONJUNCTION,

			"Иначе ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"sk": &Dialect{
		Language: "sk",
		Name:     "Slovak",
		Native:   "Slovensky",
		Keywords: map[string][]string{
			feature: {
				"Požiadavka",
				"Funkcia",
				"Vlastnosť",
			},
			rule: {
				"Rule",
			},
			background: {
				"Pozadie",
			},
			scenario: {
				"Príklad",
				"Scenár",
			},
			scenarioOutline: {
				"Náčrt Scenáru",
				"Náčrt Scenára",
				"Osnova Scenára",
			},
			examples: {
				"Príklady",
			},
			given: {
				"* ",
				"Pokiaľ ",
				"Za predpokladu ",
			},
			when: {
				"* ",
				"Keď ",
				"Ak ",
			},
			then: {
				"* ",
				"Tak ",
				"Potom ",
			},
			and: {
				"* ",
				"A ",
				"A tiež ",
				"A taktiež ",
				"A zároveň ",
			},
			but: {
				"* ",
				"Ale ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Pokiaľ ": messages.StepKeywordType_CONTEXT,

			"Za predpokladu ": messages.StepKeywordType_CONTEXT,

			"Keď ": messages.StepKeywordType_ACTION,

			"Ak ": messages.StepKeywordType_ACTION,

			"Tak ": messages.StepKeywordType_OUTCOME,

			"Potom ": messages.StepKeywordType_OUTCOME,

			"A ": messages.StepKeywordType_CONJUNCTION,

			"A tiež ": messages.StepKeywordType_CONJUNCTION,

			"A taktiež ": messages.StepKeywordType_CONJUNCTION,

			"A zároveň ": messages.StepKeywordType_CONJUNCTION,

			"Ale ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"sl": &Dialect{
		Language: "sl",
		Name:     "Slovenian",
		Native:   "Slovenski",
		Keywords: map[string][]string{
			feature: {
				"Funkcionalnost",
				"Funkcija",
				"Možnosti",
				"Moznosti",
				"Lastnost",
				"Značilnost",
			},
			rule: {
				"Rule",
			},
			background: {
				"Kontekst",
				"Osnova",
				"Ozadje",
			},
			scenario: {
				"Primer",
				"Scenarij",
			},
			scenarioOutline: {
				"Struktura scenarija",
				"Skica",
				"Koncept",
				"Oris scenarija",
				"Osnutek",
			},
			examples: {
				"Primeri",
				"Scenariji",
			},
			given: {
				"Dano ",
				"Podano ",
				"Zaradi ",
				"Privzeto ",
			},
			when: {
				"Ko ",
				"Ce ",
				"Če ",
				"Kadar ",
			},
			then: {
				"Nato ",
				"Potem ",
				"Takrat ",
			},
			and: {
				"In ",
				"Ter ",
			},
			but: {
				"Toda ",
				"Ampak ",
				"Vendar ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Dano ": messages.StepKeywordType_CONTEXT,

			"Podano ": messages.StepKeywordType_CONTEXT,

			"Zaradi ": messages.StepKeywordType_CONTEXT,

			"Privzeto ": messages.StepKeywordType_CONTEXT,

			"Ko ": messages.StepKeywordType_ACTION,

			"Ce ": messages.StepKeywordType_ACTION,

			"Če ": messages.StepKeywordType_ACTION,

			"Kadar ": messages.StepKeywordType_ACTION,

			"Nato ": messages.StepKeywordType_OUTCOME,

			"Potem ": messages.StepKeywordType_OUTCOME,

			"Takrat ": messages.StepKeywordType_OUTCOME,

			"In ": messages.StepKeywordType_CONJUNCTION,

			"Ter ": messages.StepKeywordType_CONJUNCTION,

			"Toda ": messages.StepKeywordType_CONJUNCTION,

			"Ampak ": messages.StepKeywordType_CONJUNCTION,

			"Vendar ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"sr-Cyrl": &Dialect{
		Language: "sr-Cyrl",
		Name:     "Serbian",
		Native:   "Српски",
		Keywords: map[string][]string{
			feature: {
				"Функционалност",
				"Могућност",
				"Особина",
			},
			rule: {
				"Правило",
			},
			background: {
				"Контекст",
				"Основа",
				"Позадина",
			},
			scenario: {
				"Сценарио",
				"Пример",
			},
			scenarioOutline: {
				"Структура сценарија",
				"Скица",
				"Концепт",
			},
			examples: {
				"Примери",
				"Сценарији",
			},
			given: {
				"* ",
				"За дато ",
				"За дате ",
				"За дати ",
			},
			when: {
				"* ",
				"Када ",
				"Кад ",
			},
			then: {
				"* ",
				"Онда ",
			},
			and: {
				"* ",
				"И ",
			},
			but: {
				"* ",
				"Али ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"За дато ": messages.StepKeywordType_CONTEXT,

			"За дате ": messages.StepKeywordType_CONTEXT,

			"За дати ": messages.StepKeywordType_CONTEXT,

			"Када ": messages.StepKeywordType_ACTION,

			"Кад ": messages.StepKeywordType_ACTION,

			"Онда ": messages.StepKeywordType_OUTCOME,

			"И ": messages.StepKeywordType_CONJUNCTION,

			"Али ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"sr-Latn": &Dialect{
		Language: "sr-Latn",
		Name:     "Serbian (Latin)",
		Native:   "Srpski (Latinica)",
		Keywords: map[string][]string{
			feature: {
				"Funkcionalnost",
				"Mogućnost",
				"Mogucnost",
				"Osobina",
			},
			rule: {
				"Pravilo",
			},
			background: {
				"Kontekst",
				"Osnova",
				"Pozadina",
			},
			scenario: {
				"Scenario",
				"Primer",
			},
			scenarioOutline: {
				"Struktura scenarija",
				"Skica",
				"Koncept",
			},
			examples: {
				"Primeri",
				"Scenariji",
			},
			given: {
				"* ",
				"Za dato ",
				"Za date ",
				"Za dati ",
			},
			when: {
				"* ",
				"Kada ",
				"Kad ",
			},
			then: {
				"* ",
				"Onda ",
			},
			and: {
				"* ",
				"I ",
			},
			but: {
				"* ",
				"Ali ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Za dato ": messages.StepKeywordType_CONTEXT,

			"Za date ": messages.StepKeywordType_CONTEXT,

			"Za dati ": messages.StepKeywordType_CONTEXT,

			"Kada ": messages.StepKeywordType_ACTION,

			"Kad ": messages.StepKeywordType_ACTION,

			"Onda ": messages.StepKeywordType_OUTCOME,

			"I ": messages.StepKeywordType_CONJUNCTION,

			"Ali ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"sv": &Dialect{
		Language: "sv",
		Name:     "Swedish",
		Native:   "Svenska",
		Keywords: map[string][]string{
			feature: {
				"Egenskap",
			},
			rule: {
				"Regel",
			},
			background: {
				"Bakgrund",
			},
			scenario: {
				"Scenario",
			},
			scenarioOutline: {
				"Abstrakt Scenario",
				"Scenariomall",
			},
			examples: {
				"Exempel",
			},
			given: {
				"* ",
				"Givet ",
			},
			when: {
				"* ",
				"När ",
			},
			then: {
				"* ",
				"Så ",
			},
			and: {
				"* ",
				"Och ",
			},
			but: {
				"* ",
				"Men ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Givet ": messages.StepKeywordType_CONTEXT,

			"När ": messages.StepKeywordType_ACTION,

			"Så ": messages.StepKeywordType_OUTCOME,

			"Och ": messages.StepKeywordType_CONJUNCTION,

			"Men ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ta": &Dialect{
		Language: "ta",
		Name:     "Tamil",
		Native:   "தமிழ்",
		Keywords: map[string][]string{
			feature: {
				"அம்சம்",
				"வணிக தேவை",
				"திறன்",
			},
			rule: {
				"Rule",
			},
			background: {
				"பின்னணி",
			},
			scenario: {
				"உதாரணமாக",
				"காட்சி",
			},
			scenarioOutline: {
				"காட்சி சுருக்கம்",
				"காட்சி வார்ப்புரு",
			},
			examples: {
				"எடுத்துக்காட்டுகள்",
				"காட்சிகள்",
				"நிலைமைகளில்",
			},
			given: {
				"* ",
				"கொடுக்கப்பட்ட ",
			},
			when: {
				"* ",
				"எப்போது ",
			},
			then: {
				"* ",
				"அப்பொழுது ",
			},
			and: {
				"* ",
				"மேலும் ",
				"மற்றும் ",
			},
			but: {
				"* ",
				"ஆனால் ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"கொடுக்கப்பட்ட ": messages.StepKeywordType_CONTEXT,

			"எப்போது ": messages.StepKeywordType_ACTION,

			"அப்பொழுது ": messages.StepKeywordType_OUTCOME,

			"மேலும் ": messages.StepKeywordType_CONJUNCTION,

			"மற்றும் ": messages.StepKeywordType_CONJUNCTION,

			"ஆனால் ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"th": &Dialect{
		Language: "th",
		Name:     "Thai",
		Native:   "ไทย",
		Keywords: map[string][]string{
			feature: {
				"โครงหลัก",
				"ความต้องการทางธุรกิจ",
				"ความสามารถ",
			},
			rule: {
				"Rule",
			},
			background: {
				"แนวคิด",
			},
			scenario: {
				"เหตุการณ์",
			},
			scenarioOutline: {
				"สรุปเหตุการณ์",
				"โครงสร้างของเหตุการณ์",
			},
			examples: {
				"ชุดของตัวอย่าง",
				"ชุดของเหตุการณ์",
			},
			given: {
				"* ",
				"กำหนดให้ ",
			},
			when: {
				"* ",
				"เมื่อ ",
			},
			then: {
				"* ",
				"ดังนั้น ",
			},
			and: {
				"* ",
				"และ ",
			},
			but: {
				"* ",
				"แต่ ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"กำหนดให้ ": messages.StepKeywordType_CONTEXT,

			"เมื่อ ": messages.StepKeywordType_ACTION,

			"ดังนั้น ": messages.StepKeywordType_OUTCOME,

			"และ ": messages.StepKeywordType_CONJUNCTION,

			"แต่ ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"te": &Dialect{
		Language: "te",
		Name:     "Telugu",
		Native:   "తెలుగు",
		Keywords: map[string][]string{
			feature: {
				"గుణము",
			},
			rule: {
				"Rule",
			},
			background: {
				"నేపథ్యం",
			},
			scenario: {
				"ఉదాహరణ",
				"సన్నివేశం",
			},
			scenarioOutline: {
				"కథనం",
			},
			examples: {
				"ఉదాహరణలు",
			},
			given: {
				"* ",
				"చెప్పబడినది ",
			},
			when: {
				"* ",
				"ఈ పరిస్థితిలో ",
			},
			then: {
				"* ",
				"అప్పుడు ",
			},
			and: {
				"* ",
				"మరియు ",
			},
			but: {
				"* ",
				"కాని ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"చెప్పబడినది ": messages.StepKeywordType_CONTEXT,

			"ఈ పరిస్థితిలో ": messages.StepKeywordType_ACTION,

			"అప్పుడు ": messages.StepKeywordType_OUTCOME,

			"మరియు ": messages.StepKeywordType_CONJUNCTION,

			"కాని ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"tlh": &Dialect{
		Language: "tlh",
		Name:     "Klingon",
		Native:   "tlhIngan",
		Keywords: map[string][]string{
			feature: {
				"Qap",
				"Qu'meH 'ut",
				"perbogh",
				"poQbogh malja'",
				"laH",
			},
			rule: {
				"Rule",
			},
			background: {
				"mo'",
			},
			scenario: {
				"lut",
			},
			scenarioOutline: {
				"lut chovnatlh",
			},
			examples: {
				"ghantoH",
				"lutmey",
			},
			given: {
				"* ",
				"ghu' noblu' ",
				"DaH ghu' bejlu' ",
			},
			when: {
				"* ",
				"qaSDI' ",
			},
			then: {
				"* ",
				"vaj ",
			},
			and: {
				"* ",
				"'ej ",
				"latlh ",
			},
			but: {
				"* ",
				"'ach ",
				"'a ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"ghu' noblu' ": messages.StepKeywordType_CONTEXT,

			"DaH ghu' bejlu' ": messages.StepKeywordType_CONTEXT,

			"qaSDI' ": messages.StepKeywordType_ACTION,

			"vaj ": messages.StepKeywordType_OUTCOME,

			"'ej ": messages.StepKeywordType_CONJUNCTION,

			"latlh ": messages.StepKeywordType_CONJUNCTION,

			"'ach ": messages.StepKeywordType_CONJUNCTION,

			"'a ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"tr": &Dialect{
		Language: "tr",
		Name:     "Turkish",
		Native:   "Türkçe",
		Keywords: map[string][]string{
			feature: {
				"Özellik",
				"İş Gereksinimi",
				"Gereksinim",
				"İşlev",
				"Kullanıcı Hikayesi",
				"Yetenek",
				"Teknik Gereksinim",
			},
			rule: {
				"Kural",
				"İş Kuralı",
				"Kaide",
				"Hüküm",
				"Madde",
			},
			background: {
				"Geçmiş",
				"Arka Plan",
				"Ön Koşul",
				"Önkoşul",
				"Önceki Durum",
				"Giriş",
				"Mukaddime",
				"Mevcut Durum",
			},
			scenario: {
				"Örnek",
				"Senaryo",
				"Durum",
				"Vaka",
			},
			scenarioOutline: {
				"Senaryo taslağı",
				"Senaryo şablonu",
			},
			examples: {
				"Örnekler",
				"Değerler",
			},
			given: {
				"* ",
				"Mevcut ",
				"Önceden ",
				"Geçmişte ",
				"Daha önce ",
				"Halihazırda ",
				"Zaten ",
				"Sistemde ",
				"Diyelim ki ",
				"Varsayalım ki ",
				"Farz edelim ki ",
				"Kabul edelim ki ",
				"Başlangıçta ",
				"Varsayılan olarak ",
				"Biliniyor ki ",
			},
			when: {
				"* ",
				"Eğer ",
				"Eğer ki ",
				"Ne zaman ",
				"Ne zaman ki ",
				"Şayet ",
			},
			then: {
				"* ",
				"Beklenen ",
				"O zaman ",
				"Sonuç olarak ",
				"Böylece ",
				"Bunun üzerine ",
				"Bu durumda ",
				"O takdirde ",
				"Şu halde ",
				"Netice itibariyle ",
				"Buna binaen ",
			},
			and: {
				"* ",
				"Ve ",
				"Hem de ",
				"Bir de ",
				"Ayrıca ",
				"İlaveten ",
				"Buna ek olarak ",
			},
			but: {
				"* ",
				"Fakat ",
				"Ama ",
				"Ancak ",
				"Yalnız ",
				"Lakin ",
				"Meğer ki ",
				"Buna mukabil ",
				"Aksi halde ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Mevcut ": messages.StepKeywordType_CONTEXT,

			"Önceden ": messages.StepKeywordType_CONTEXT,

			"Geçmişte ": messages.StepKeywordType_CONTEXT,

			"Daha önce ": messages.StepKeywordType_CONTEXT,

			"Halihazırda ": messages.StepKeywordType_CONTEXT,

			"Zaten ": messages.StepKeywordType_CONTEXT,

			"Sistemde ": messages.StepKeywordType_CONTEXT,

			"Diyelim ki ": messages.StepKeywordType_CONTEXT,

			"Varsayalım ki ": messages.StepKeywordType_CONTEXT,

			"Farz edelim ki ": messages.StepKeywordType_CONTEXT,

			"Kabul edelim ki ": messages.StepKeywordType_CONTEXT,

			"Başlangıçta ": messages.StepKeywordType_CONTEXT,

			"Varsayılan olarak ": messages.StepKeywordType_CONTEXT,

			"Biliniyor ki ": messages.StepKeywordType_CONTEXT,

			"Eğer ": messages.StepKeywordType_ACTION,

			"Eğer ki ": messages.StepKeywordType_ACTION,

			"Ne zaman ": messages.StepKeywordType_ACTION,

			"Ne zaman ki ": messages.StepKeywordType_ACTION,

			"Şayet ": messages.StepKeywordType_ACTION,

			"Beklenen ": messages.StepKeywordType_OUTCOME,

			"O zaman ": messages.StepKeywordType_OUTCOME,

			"Sonuç olarak ": messages.StepKeywordType_OUTCOME,

			"Böylece ": messages.StepKeywordType_OUTCOME,

			"Bunun üzerine ": messages.StepKeywordType_OUTCOME,

			"Bu durumda ": messages.StepKeywordType_OUTCOME,

			"O takdirde ": messages.StepKeywordType_OUTCOME,

			"Şu halde ": messages.StepKeywordType_OUTCOME,

			"Netice itibariyle ": messages.StepKeywordType_OUTCOME,

			"Buna binaen ": messages.StepKeywordType_OUTCOME,

			"Ve ": messages.StepKeywordType_CONJUNCTION,

			"Hem de ": messages.StepKeywordType_CONJUNCTION,

			"Bir de ": messages.StepKeywordType_CONJUNCTION,

			"Ayrıca ": messages.StepKeywordType_CONJUNCTION,

			"İlaveten ": messages.StepKeywordType_CONJUNCTION,

			"Buna ek olarak ": messages.StepKeywordType_CONJUNCTION,

			"Fakat ": messages.StepKeywordType_CONJUNCTION,

			"Ama ": messages.StepKeywordType_CONJUNCTION,

			"Ancak ": messages.StepKeywordType_CONJUNCTION,

			"Yalnız ": messages.StepKeywordType_CONJUNCTION,

			"Lakin ": messages.StepKeywordType_CONJUNCTION,

			"Meğer ki ": messages.StepKeywordType_CONJUNCTION,

			"Buna mukabil ": messages.StepKeywordType_CONJUNCTION,

			"Aksi halde ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"tt": &Dialect{
		Language: "tt",
		Name:     "Tatar",
		Native:   "Татарча",
		Keywords: map[string][]string{
			feature: {
				"Мөмкинлек",
				"Үзенчәлеклелек",
			},
			rule: {
				"Rule",
			},
			background: {
				"Кереш",
			},
			scenario: {
				"Сценарий",
			},
			scenarioOutline: {
				"Сценарийның төзелеше",
			},
			examples: {
				"Үрнәкләр",
				"Мисаллар",
			},
			given: {
				"* ",
				"Әйтик ",
			},
			when: {
				"* ",
				"Әгәр ",
			},
			then: {
				"* ",
				"Нәтиҗәдә ",
			},
			and: {
				"* ",
				"Һәм ",
				"Вә ",
			},
			but: {
				"* ",
				"Ләкин ",
				"Әмма ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Әйтик ": messages.StepKeywordType_CONTEXT,

			"Әгәр ": messages.StepKeywordType_ACTION,

			"Нәтиҗәдә ": messages.StepKeywordType_OUTCOME,

			"Һәм ": messages.StepKeywordType_CONJUNCTION,

			"Вә ": messages.StepKeywordType_CONJUNCTION,

			"Ләкин ": messages.StepKeywordType_CONJUNCTION,

			"Әмма ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"uk": &Dialect{
		Language: "uk",
		Name:     "Ukrainian",
		Native:   "Українська",
		Keywords: map[string][]string{
			feature: {
				"Функціонал",
			},
			rule: {
				"Rule",
			},
			background: {
				"Передумова",
			},
			scenario: {
				"Приклад",
				"Сценарій",
			},
			scenarioOutline: {
				"Структура сценарію",
			},
			examples: {
				"Приклади",
			},
			given: {
				"* ",
				"Припустимо ",
				"Припустимо, що ",
				"Нехай ",
				"Дано ",
			},
			when: {
				"* ",
				"Якщо ",
				"Коли ",
			},
			then: {
				"* ",
				"То ",
				"Тоді ",
			},
			and: {
				"* ",
				"І ",
				"А також ",
				"Та ",
			},
			but: {
				"* ",
				"Але ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Припустимо ": messages.StepKeywordType_CONTEXT,

			"Припустимо, що ": messages.StepKeywordType_CONTEXT,

			"Нехай ": messages.StepKeywordType_CONTEXT,

			"Дано ": messages.StepKeywordType_CONTEXT,

			"Якщо ": messages.StepKeywordType_ACTION,

			"Коли ": messages.StepKeywordType_ACTION,

			"То ": messages.StepKeywordType_OUTCOME,

			"Тоді ": messages.StepKeywordType_OUTCOME,

			"І ": messages.StepKeywordType_CONJUNCTION,

			"А також ": messages.StepKeywordType_CONJUNCTION,

			"Та ": messages.StepKeywordType_CONJUNCTION,

			"Але ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ur": &Dialect{
		Language: "ur",
		Name:     "Urdu",
		Native:   "اردو",
		Keywords: map[string][]string{
			feature: {
				"صلاحیت",
				"کاروبار کی ضرورت",
				"خصوصیت",
			},
			rule: {
				"Rule",
			},
			background: {
				"پس منظر",
			},
			scenario: {
				"منظرنامہ",
			},
			scenarioOutline: {
				"منظر نامے کا خاکہ",
			},
			examples: {
				"مثالیں",
			},
			given: {
				"* ",
				"اگر ",
				"بالفرض ",
				"فرض کیا ",
			},
			when: {
				"* ",
				"جب ",
			},
			then: {
				"* ",
				"پھر ",
				"تب ",
			},
			and: {
				"* ",
				"اور ",
			},
			but: {
				"* ",
				"لیکن ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"اگر ": messages.StepKeywordType_CONTEXT,

			"بالفرض ": messages.StepKeywordType_CONTEXT,

			"فرض کیا ": messages.StepKeywordType_CONTEXT,

			"جب ": messages.StepKeywordType_ACTION,

			"پھر ": messages.StepKeywordType_OUTCOME,

			"تب ": messages.StepKeywordType_OUTCOME,

			"اور ": messages.StepKeywordType_CONJUNCTION,

			"لیکن ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"uz": &Dialect{
		Language: "uz",
		Name:     "Uzbek",
		Native:   "Узбекча",
		Keywords: map[string][]string{
			feature: {
				"Функционал",
			},
			rule: {
				"Rule",
			},
			background: {
				"Тарих",
			},
			scenario: {
				"Сценарий",
			},
			scenarioOutline: {
				"Сценарий структураси",
			},
			examples: {
				"Мисоллар",
			},
			given: {
				"* ",
				"Belgilangan ",
			},
			when: {
				"* ",
				"Агар ",
			},
			then: {
				"* ",
				"Унда ",
			},
			and: {
				"* ",
				"Ва ",
			},
			but: {
				"* ",
				"Лекин ",
				"Бирок ",
				"Аммо ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Belgilangan ": messages.StepKeywordType_CONTEXT,

			"Агар ": messages.StepKeywordType_ACTION,

			"Унда ": messages.StepKeywordType_OUTCOME,

			"Ва ": messages.StepKeywordType_CONJUNCTION,

			"Лекин ": messages.StepKeywordType_CONJUNCTION,

			"Бирок ": messages.StepKeywordType_CONJUNCTION,

			"Аммо ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"vi": &Dialect{
		Language: "vi",
		Name:     "Vietnamese",
		Native:   "Tiếng Việt",
		Keywords: map[string][]string{
			feature: {
				"Tính năng",
			},
			rule: {
				"Quy tắc",
			},
			background: {
				"Bối cảnh",
			},
			scenario: {
				"Tình huống",
				"Kịch bản",
			},
			scenarioOutline: {
				"Khung tình huống",
				"Khung kịch bản",
			},
			examples: {
				"Dữ liệu",
			},
			given: {
				"* ",
				"Biết ",
				"Cho ",
			},
			when: {
				"* ",
				"Khi ",
			},
			then: {
				"* ",
				"Thì ",
			},
			and: {
				"* ",
				"Và ",
			},
			but: {
				"* ",
				"Nhưng ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"Biết ": messages.StepKeywordType_CONTEXT,

			"Cho ": messages.StepKeywordType_CONTEXT,

			"Khi ": messages.StepKeywordType_ACTION,

			"Thì ": messages.StepKeywordType_OUTCOME,

			"Và ": messages.StepKeywordType_CONJUNCTION,

			"Nhưng ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"zh-CN": &Dialect{
		Language: "zh-CN",
		Name:     "Chinese simplified",
		Native:   "简体中文",
		Keywords: map[string][]string{
			feature: {
				"功能",
			},
			rule: {
				"Rule",
				"规则",
			},
			background: {
				"背景",
			},
			scenario: {
				"场景",
				"剧本",
			},
			scenarioOutline: {
				"场景大纲",
				"剧本大纲",
			},
			examples: {
				"例子",
			},
			given: {
				"* ",
				"假如",
				"假设",
				"假定",
			},
			when: {
				"* ",
				"当",
			},
			then: {
				"* ",
				"那么",
			},
			and: {
				"* ",
				"而且",
				"并且",
				"同时",
			},
			but: {
				"* ",
				"但是",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"假如": messages.StepKeywordType_CONTEXT,

			"假设": messages.StepKeywordType_CONTEXT,

			"假定": messages.StepKeywordType_CONTEXT,

			"当": messages.StepKeywordType_ACTION,

			"那么": messages.StepKeywordType_OUTCOME,

			"而且": messages.StepKeywordType_CONJUNCTION,

			"并且": messages.StepKeywordType_CONJUNCTION,

			"同时": messages.StepKeywordType_CONJUNCTION,

			"但是": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"ml": &Dialect{
		Language: "ml",
		Name:     "Malayalam",
		Native:   "മലയാളം",
		Keywords: map[string][]string{
			feature: {
				"സവിശേഷത",
			},
			rule: {
				"നിയമം",
			},
			background: {
				"പശ്ചാത്തലം",
			},
			scenario: {
				"രംഗം",
			},
			scenarioOutline: {
				"സാഹചര്യത്തിന്റെ രൂപരേഖ",
			},
			examples: {
				"ഉദാഹരണങ്ങൾ",
			},
			given: {
				"* ",
				"നൽകിയത്",
			},
			when: {
				"എപ്പോൾ",
			},
			then: {
				"* ",
				"പിന്നെ",
			},
			and: {
				"* ",
				"ഒപ്പം",
			},
			but: {
				"* ",
				"പക്ഷേ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"നൽകിയത്": messages.StepKeywordType_CONTEXT,

			"എപ്പോൾ": messages.StepKeywordType_ACTION,

			"പിന്നെ": messages.StepKeywordType_OUTCOME,

			"ഒപ്പം": messages.StepKeywordType_CONJUNCTION,

			"പക്ഷേ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"zh-TW": &Dialect{
		Language: "zh-TW",
		Name:     "Chinese traditional",
		Native:   "繁體中文",
		Keywords: map[string][]string{
			feature: {
				"功能",
			},
			rule: {
				"Rule",
			},
			background: {
				"背景",
			},
			scenario: {
				"場景",
				"劇本",
			},
			scenarioOutline: {
				"場景大綱",
				"劇本大綱",
			},
			examples: {
				"例子",
			},
			given: {
				"* ",
				"假如",
				"假設",
				"假定",
			},
			when: {
				"* ",
				"當",
			},
			then: {
				"* ",
				"那麼",
			},
			and: {
				"* ",
				"而且",
				"並且",
				"同時",
			},
			but: {
				"* ",
				"但是",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"假如": messages.StepKeywordType_CONTEXT,

			"假設": messages.StepKeywordType_CONTEXT,

			"假定": messages.StepKeywordType_CONTEXT,

			"當": messages.StepKeywordType_ACTION,

			"那麼": messages.StepKeywordType_OUTCOME,

			"而且": messages.StepKeywordType_CONJUNCTION,

			"並且": messages.StepKeywordType_CONJUNCTION,

			"同時": messages.StepKeywordType_CONJUNCTION,

			"但是": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"mr": &Dialect{
		Language: "mr",
		Name:     "Marathi",
		Native:   "मराठी",
		Keywords: map[string][]string{
			feature: {
				"वैशिष्ट्य",
				"सुविधा",
			},
			rule: {
				"नियम",
			},
			background: {
				"पार्श्वभूमी",
			},
			scenario: {
				"परिदृश्य",
			},
			scenarioOutline: {
				"परिदृश्य रूपरेखा",
			},
			examples: {
				"उदाहरण",
			},
			given: {
				"* ",
				"जर",
				"दिलेल्या प्रमाणे ",
			},
			when: {
				"* ",
				"जेव्हा ",
			},
			then: {
				"* ",
				"मग ",
				"तेव्हा ",
			},
			and: {
				"* ",
				"आणि ",
				"तसेच ",
			},
			but: {
				"* ",
				"पण ",
				"परंतु ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"जर": messages.StepKeywordType_CONTEXT,

			"दिलेल्या प्रमाणे ": messages.StepKeywordType_CONTEXT,

			"जेव्हा ": messages.StepKeywordType_ACTION,

			"मग ": messages.StepKeywordType_OUTCOME,

			"तेव्हा ": messages.StepKeywordType_OUTCOME,

			"आणि ": messages.StepKeywordType_CONJUNCTION,

			"तसेच ": messages.StepKeywordType_CONJUNCTION,

			"पण ": messages.StepKeywordType_CONJUNCTION,

			"परंतु ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
	"amh": &Dialect{
		Language: "amh",
		Name:     "Amharic",
		Native:   "አማርኛ",
		Keywords: map[string][]string{
			feature: {
				"ስራ",
				"የተፈለገው ስራ",
				"የሚፈለገው ድርጊት",
			},
			rule: {
				"ህግ",
			},
			background: {
				"ቅድመ ሁኔታ",
				"መነሻ",
				"መነሻ ሀሳብ",
			},
			scenario: {
				"ምሳሌ",
				"ሁናቴ",
			},
			scenarioOutline: {
				"ሁናቴ ዝርዝር",
				"ሁናቴ አብነት",
			},
			examples: {
				"ምሳሌዎች",
				"ሁናቴዎች",
			},
			given: {
				"* ",
				"የተሰጠ ",
			},
			when: {
				"* ",
				"መቼ ",
			},
			then: {
				"* ",
				"ከዚያ ",
			},
			and: {
				"* ",
				"እና ",
			},
			but: {
				"* ",
				"ግን ",
			},
		},
		KeywordTypes: map[string]messages.StepKeywordType{
			"የተሰጠ ": messages.StepKeywordType_CONTEXT,

			"መቼ ": messages.StepKeywordType_ACTION,

			"ከዚያ ": messages.StepKeywordType_OUTCOME,

			"እና ": messages.StepKeywordType_CONJUNCTION,

			"ግን ": messages.StepKeywordType_CONJUNCTION,

			"* ": messages.StepKeywordType_UNKNOWN,
		},
	},
}
