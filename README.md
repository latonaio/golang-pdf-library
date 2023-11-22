# README #

## TODO
- apply godoc
- path ../

## Specifications

- data.jsonの値は文字列のみ(数値は不可、フォーマットした文字列)
- AutoScale is not applicable to multiline text 

## What is this repository for? ##

以下をソースとしてpdfを生成する。
- template.json 
- data.json 
- template.pdf (template.jsonにて指定)

※各ファイル名はパラメータで変更可能

## Parameters ##
サンプル実行用パラメータ
- -g draw grid (default fasle)
- -b fill background (default false)
- -s fill sample text (default false)
- -t template json file name (default "./resources/sampleTemplate.json")
- -d data source json file name (default "./resources/sampleData.json")
- -p path of private key for signing (if need signature)
- -c path of certification for signing (if need signature)
- -h path of chain for signing (if need chain)
- -o generate pdf file name (default "sample.pdf")

## Formats ##

### template json ###

```json:template.json
s{
    "version" : "1.0",                              // pdf module version
    "orientation" : "landscape",                    // [landscape|portlait]
    "size" : "A4",                                  // size
    "image" : "millcert.pdf",                       // background image pdf
    "styles" : {                                    // define styles 
        "default" : {                               // default style
            "fontFamily" : "xxxxx",             
            "fontSize" : 9,
            "color" : "#000",
            "underline" : false,
            "horizontalAlignment" : "left",
            "verticalAlignment" : "middle"
        },
        "default-numeric" : {
            "style" : "default",                    // base style
            "horizontalAlignment" : "right"
        },
        ...
    },
    "fields" : [                                    // define fields
        {                                           // simple field
            "dataSource" : "customer",              // data source(this must be value in data.json)
            "style" : "default",                    // style(from styles)
            "rect" : [                              // rect to draw field
                180, 100, 220, 12                   // [x, y, width, height]
            ]
        },
        ...
        {                                           // record field
            "dataSource" : "properties",            // data source(this must be array in data.json)
            "rect" : [                              // rect to draw records
                140, 293, 263, 90
            ],
            "record" : {                            // define as record
                "direction" : "x",                  // record direction
                "size" : [                          // size to draw record
                    52.6, 90                        // [x, y]
                ],
                "fields" : [                        // define fields on record
                    {
                        "dataSource" : "labelEn",   // data source(this must be records' member in data.json)
                        "style" : "default",        // style(from styles)
                        "rect" : [                  // offset rect in record
                            0, 0, 52.6, 15
                        ]
                    },
                    {
                        "dataSource" : "max",
                        "style" : "default",
                        "rect" : [
                            0, 15, 52.6, 15         // offset rect in record
                        ]
                    },
                    {
                        "dataSource" : "unit",
                        "style" : "default",
                        "rect" : [
                            0, 30, 52.6, 15         // offset rect in record
                        ]
                    },
                    ...
                ]
            }
        },
        ...
    ]
}
```

### data json ###

```json:data.json
{
    "customer" : "Mill Cert Sample Customer",
    ...
    "properties" : [
        {
            "labelEn" : "Yield Strength",
            "max" : "245.000",
            "unit" : "N/m ㎡"
        },
        {
            "labelEn" : "Tensile Strength",
            "max" : "640.500",
            "unit" : "N/m ㎡"
        },
        ...
    ]
}```