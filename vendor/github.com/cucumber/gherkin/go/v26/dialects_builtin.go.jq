. as $root
| (
  [ to_entries[]
    | [
        "\t",(.key|@json),": &Dialect{\n",
        "\t\t", (.key|@json),", ", (.value.name|@json),", ", (.value.native|@json), ", map[string][]string{\n"
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
        "\t\tmap[string]messages.StepKeywordType{\n"
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
        "\t\t}",
        "},\n"
      ]
    | add
  ]
  | add
  )
| "package gherkin\n\n"
+ "import messages \"github.com/cucumber/messages/go/v21\"\n\n"
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
