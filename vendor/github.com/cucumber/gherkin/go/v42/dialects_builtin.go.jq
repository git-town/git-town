. as $root
| (
  [ to_entries[]
    | [
        "\t",(.key|@json),": &Dialect{\n",
        "\t\tLanguage: ", (.key|@json), ",\n",
        "\t\tName:     ", (.value.name|@json),",\n",
        "\t\tNative:   ",(.value.native|@json),",\n",
        "\t\tKeywords: map[string][]string{\n"
      ] + (
          [ .value
            | {"feature","rule","background","scenario","scenarioOutline","examples","given","when","then","and","but"}
            | to_entries[]
            | "\t\t\t"+(.key), ": {\n",
              ([ .value[] | "\t\t\t\t", @json, ",\n"  ]|add),
              "\t\t\t},\n"
          ]
      ) + [
        "\t\t},\n",
        "\t\tKeywordTypes: map[string]messages.StepKeywordType{\n"
      ] + (
        [ .value.given
          | (
              [ .[] | select(. != "* ") |
                "\t\t\t",
                @json,
                ": messages.StepKeywordType_CONTEXT",
                ",\n\n"
              ] | add
            ),
            ""
        ]
        +
        [ .value.when
          | (
              [ .[] | select(. != "* ") |
                "\t\t\t",
                @json,
                ": messages.StepKeywordType_ACTION",
                ",\n\n"
              ] | add
            ),
            ""
        ]
        +
        [ .value.then
          | (
              [ .[] | select(. != "* ") |
                "\t\t\t",
                @json,
                ": messages.StepKeywordType_OUTCOME",
                ",\n\n"
              ] | add
            ),
            ""
        ]
        +
        [ .value.and
          | (
              [ .[] | select(. != "* ") |
                "\t\t\t",
                @json,
                ": messages.StepKeywordType_CONJUNCTION",
                ",\n\n"
              ] | add
            ),
            ""
        ]
        +
        [ .value.but
          | (
              [ .[] | select(. != "* ") |
                "\t\t\t",
                @json,
                ": messages.StepKeywordType_CONJUNCTION",
                ",\n\n"
              ] | add
            ),
            ""
        ]
        + [
          "\t\t\t\"* \": messages.StepKeywordType_UNKNOWN,\n"
        ]
      ) + [
        "\t\t},\n",
        "\t},\n"
      ]
    | add
  ]
  | add
  )
| "// Code generated from dialects_builtin.go.jq (make dialects_builtin.go); DO NOT EDIT.\n\n" # Standard header defined at https://golang.org/s/generatedcode
+ "package gherkin\n\n"
+ "import messages \"github.com/cucumber/messages/go/v34\"\n\n"
+ "// Builtin dialects for " + ([ $root | to_entries[] | .key+" ("+.value.name+")" ] | join(", ")) + "\n"
+ "func DialectsBuiltin() DialectProvider {\n"
+ "\treturn builtinDialects\n"
+ "}\n\n"
+ "const (\n"
+ "	feature         = \"feature\"\n"
+ "	rule            = \"rule\"\n"
+ "	background      = \"background\"\n"
+ "	scenario        = \"scenario\"\n"
+ "	scenarioOutline = \"scenarioOutline\"\n"
+ "	examples        = \"examples\"\n"
+ "	given           = \"given\"\n"
+ "	when            = \"when\"\n"
+ "	then            = \"then\"\n"
+ "	and             = \"and\"\n"
+ "	but             = \"but\"\n"
+ ")\n\n"
+ "var builtinDialects = gherkinDialectMap{\n"
+ .
+ "}"
