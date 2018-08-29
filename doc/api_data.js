define({ "api": [
  {
    "type": "post",
    "url": "/create",
    "title": "To create Multisig wallet.",
    "name": "Create_Multisig_wallet",
    "group": "MultisigWallet",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "txbytes",
            "description": "<p>Transaction bytes to initiate transaction for wallet, by default its empty at first initiation.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Name of Account.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>Password for account.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "chain_id",
            "description": "<p>Chain Id.</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "account_number",
            "description": "<p>Account number.</p>"
          },
          {
            "group": "Parameter",
            "type": "number",
            "optional": false,
            "field": "gas",
            "description": "<p>Gas value.</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "AccountAlreadyExists",
            "description": "<p>AccountName is  already exists</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "AccountAlreadyExists-Response:",
          "content": "{\n  Account with name XXXXX... already exists.\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Response:",
          "content": "{\n \"check_tx\": {\n   \"log\": \"Msg 0: \",\n   \"gasWanted\": \"21000\",\n   \"gasUsed\": \"1209\"\n },\n \"deliver_tx\": {\n   \"data\": \"IGpXqASE+6AvgVVRO3NyNtCzqH4=\",\n   \"log\": \"Msg 0: \",\n   \"gasWanted\": \"21000\",\n   \"gasUsed\": \"6670\",\n   \"tags\": [\n     {\n       \"key\": \"bXVsdGlzaWcgYWRkZHJlc3M=\",\n       \"value\": \"Y29zbW9zYWNjYWRkcjF5cDQ5MDJxeXNuYTZxdHVwMjRnbmt1bWp4bWd0ODJyN2g1a3V2Yw==\"\n     }\n   ]\n },\n \"hash\": \"CC78A0E5445A2EE945308F3A599EF96BD529A9AF\",\n \"height\": \"14863\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "modules/multisig/client/rest/create.go",
    "groupTitle": "MultisigWallet"
  },
  {
    "type": "post",
    "url": "/transfer",
    "title": "To send tokens from Multisig wallet.",
    "name": "Transfer_tokens_from_Multisig_wallet",
    "group": "MultisigWallet",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "to",
            "description": "<p>To Address.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Name of Account.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>Password for account.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "chain_id",
            "description": "<p>Chain Id.</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "account_number",
            "description": "<p>Account number.</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "gas",
            "description": "<p>Gas value.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "amount",
            "description": "<p>amount to send.</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "AccountAlreadyExists",
            "description": "<p>AccountName is  already exists</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "AccountAlreadyExists-Response:",
          "content": "{\n  Account with name XXXXX... already exists.\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Response:",
          "content": "{\n \"check_tx\": {\n   \"log\": \"Msg 0: \",\n   \"gasWanted\": \"21000\",\n   \"gasUsed\": \"1209\"\n },\n \"deliver_tx\": {\n   \"data\": \"IGpXqASE+6AvgVVRO3NyNtCzqH4=\",\n   \"log\": \"Msg 0: \",\n   \"gasWanted\": \"21000\",\n   \"gasUsed\": \"6670\",\n   \"tags\": [\n     {\n       \"key\": \"bXVsdGlzaWcgYWRkZHJlc3M=\",\n       \"value\": \"Y29zbW9zYWNjYWRkcjF5cDQ5MDJxeXNuYTZxdHVwMjRnbmt1bWp4bWd0ODJyN2g1a3V2Yw==\"\n     }\n   ]\n },\n \"hash\": \"CC78A0E5445A2EE945308F3A599EF96BD529A9AF\",\n \"height\": \"14863\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "modules/multisig/client/rest/send.go",
    "groupTitle": "MultisigWallet"
  },
  {
    "type": "post",
    "url": "/send",
    "title": "To send tokens to Multisig wallet.",
    "name": "Transfer_tokens_to_Multisig_wallet",
    "group": "MultisigWallet",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "txbytes",
            "description": "<p>Transaction bytes.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Name of Account.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>Password for account.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "chain_id",
            "description": "<p>Chain Id.</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "account_number",
            "description": "<p>Account number.</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "gas",
            "description": "<p>Gas value.</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "AccountAlreadyExists",
            "description": "<p>AccountName is  already exists</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "AccountAlreadyExists-Response:",
          "content": "{\n  Account with name XXXXX... already exists.\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Response:",
          "content": "{\n \"check_tx\": {\n   \"log\": \"Msg 0: \",\n   \"gasWanted\": \"21000\",\n   \"gasUsed\": \"1209\"\n },\n \"deliver_tx\": {\n   \"data\": \"IGpXqASE+6AvgVVRO3NyNtCzqH4=\",\n   \"log\": \"Msg 0: \",\n   \"gasWanted\": \"21000\",\n   \"gasUsed\": \"6670\",\n   \"tags\": [\n     {\n       \"key\": \"bXVsdGlzaWcgYWRkZHJlc3M=\",\n       \"value\": \"Y29zbW9zYWNjYWRkcjF5cDQ5MDJxeXNuYTZxdHVwMjRnbmt1bWp4bWd0ODJyN2g1a3V2Yw==\"\n     }\n   ]\n },\n \"hash\": \"CC78A0E5445A2EE945308F3A599EF96BD529A9AF\",\n \"height\": \"14863\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "modules/multisig/client/rest/fund.go",
    "groupTitle": "MultisigWallet"
  },
  {
    "type": "post",
    "url": "/initiate",
    "title": "To initiate transaction for Multisig wallet.",
    "name": "initiate_Multisig_wallet",
    "group": "MultisigWallet",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "txbytes",
            "description": "<p>Transaction bytes to initiate transaction for wallet, by default its empty at first initiation.</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "min_keys",
            "description": "<p>Number of minimum signatures required.</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "total_keys",
            "description": "<p>Number of total signatures including.</p>"
          },
          {
            "group": "Parameter",
            "type": "Boolean",
            "optional": false,
            "field": "order",
            "description": "<p>Order of signatures required or not.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Name of Account.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>Password for account.</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "AccountAlreadyExists",
            "description": "<p>AccountName is  already exists</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "AccountAlreadyExists-Response:",
          "content": "{\n  Account with name XXXXX... already exists.\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Response:",
          "content": "{\n vgEIBBAEGAEibFkyOXpiVzl6WVdOamNIVmlNV0ZrWkhkdWNHVndjWFJsTldVeWNqVmhkbmgwTWpaa2NUbDVjWEowZERSeE9UQ\n jVZMlprYURoeGN6TXpaWFI2TldabmQzRnhkRE0xZUhCMFozYzRibVl6Y0dzPSgCMkYwRAIgSntwx54iNoDqHyYSgRdxyei2n\n EkhPs3oSWVWIcSrjgCICzsjgzA6pYMKR/w1jxJ+IUNrPasDwpMEtNi2bMdaNH2\n\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "modules/multisig/client/rest/tx.go",
    "groupTitle": "MultisigWallet"
  },
  {
    "type": "post",
    "url": "/initiate_txn",
    "title": "To initiate transaction for sending Tokens from Multisig wallet.",
    "name": "initiate_transaction_of_Multisig_wallet",
    "group": "MultisigWallet",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "txbytes",
            "description": "<p>Transaction bytes to initiate transaction for wallet, by default its empty at first initiation.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "to",
            "description": "<p>To address.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "multisig_address",
            "description": "<p>Address of the multisignature wallet.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Name of Account.</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>Password for account.</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": false,
            "field": "tx_number",
            "description": "<p>Transaction number of the multisig wallet.</p>"
          }
        ]
      }
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "AccountAlreadyExists",
            "description": "<p>AccountName is  already exists</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "AccountAlreadyExists-Response:",
          "content": "{\n  Account with name XXXXX... already exists.\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Response:",
          "content": "{\n vgEIBBAEGAEibFkyOXpiVzl6WVdOamNIVmlNV0ZrWkhkdWNHVndjWFJsTldVeWNqVmhkbmgwTWpaa2NUbDVjWEowZERSeE9UQ\n jVZMlprYURoeGN6TXpaWFI2TldabmQzRnhkRE0xZUhCMFozYzRibVl6Y0dzPSgCMkYwRAIgSntwx54iNoDqHyYSgRdxyei2n/\n EkhPs3oSWVWIcSrjgCICzsjgzA6pYMKR/w1jxJ+IUNrPasDwpMEtNi2bMdaNH2\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "modules/multisig/client/rest/spend.go",
    "groupTitle": "MultisigWallet"
  }
] });
