package registry

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	crypto2 "github.com/onflow/crypto"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-standard-transactions/template"
)

func simpleTemplateWithLoop(
	name string,
	label template.Label,
	initialLoopLength uint64,
	body string,
) *template.SimpleTemplate {
	return template.NewSimpleTemplate(
		name,
		label,
		1,
	).
		WithInitialParameters(template.Parameters{initialLoopLength}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				return template.TransactionEdit{
					PrepareBlock: template.LoopTemplate(parameters[0], body),
				}, nil
			},
		)
}

var SimpleTemplates = []template.Template{
	simpleTemplateWithLoop(
		"empty loop",
		"EL",
		5397,
		"",
	),
	simpleTemplateWithLoop(
		"Assert True",
		"Assert",
		2657,
		"assert(true)",
	),
	simpleTemplateWithLoop(
		"get signer address",
		"GSA",
		3590,
		"signer.address",
	),
	simpleTemplateWithLoop(
		"get signer public account",
		"GSAcc",
		2627,
		"getAccount(signer.address)",
	),
	simpleTemplateWithLoop(
		"get signer account balance",
		"GSAccBal",
		27,
		"signer.balance",
	),
	simpleTemplateWithLoop(
		"get signer account available balance",
		"GSAccAwBal",
		24,
		"signer.availableBalance",
	),
	simpleTemplateWithLoop(
		"get signer account storage used",
		"GSAccSU",
		670,
		"signer.storage.used",
	),
	simpleTemplateWithLoop(
		"get signer account storage capacity",
		"GSAccSC",
		26,
		"signer.storage.capacity",
	),
	simpleTemplateWithLoop(
		"borrow signer account FlowToken.Vault",
		"BFTV",
		719,
		"let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(from: /storage/flowTokenVault)!",
	),
	simpleTemplateWithLoop(
		"borrow signer account FungibleToken.Receiver",
		"BFR",
		393,
		`
			let receiverRef = getAccount(signer.address)
				.capabilities.borrow<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)!
			`,
	),
	simpleTemplateWithLoop(
		"transfer tokens to self",
		"TTS",
		26,
		`
			let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(from: /storage/flowTokenVault)!
			let receiverRef = getAccount(signer.address)
				.capabilities.borrow<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)!
			receiverRef.deposit(from: <-vaultRef.withdraw(amount: 0.00001))
			`,
	),
	simpleTemplateWithLoop(
		"create new account",
		"CA",
		7,
		`let acct = Account(payer: signer)`,
	),
	simpleTemplateWithLoop(
		"create new account with contract",
		"CAWC",
		5,
		`
				let acct = Account(payer: signer)
				acct.contracts.add(name: "EmptyContract", code: "61636365737328616c6c2920636f6e747261637420456d707479436f6e7472616374207b7d".decodeHex())
			`,
	),
	simpleTemplateWithLoop(
		"decode hex",
		"HEX",
		811,
		`
				"f847b84000fb479cb398ab7e31d6f048c12ec5b5b679052589280cacde421af823f93fe927dfc3d1e371b172f97ceeac1bc235f60654184c83f4ea70dd3b7785ffb3c73802038203e8".decodeHex()
			`,
	),
	simpleTemplateWithLoop(
		"Revertible random",
		"RR",
		1958,
		`
				revertibleRandom<UInt64>(modulo: UInt64(100))
			`,
	),
	simpleTemplateWithLoop(
		"number to string conversion",
		"TS",
		2627,
		`
				i.toString()
			`,
	),
	simpleTemplateWithLoop(
		"concatenate string",
		"CS",
		2092,
		`
				"x".concat(i.toString())
			`,
	),
	// no unique intensities
	//template.NewSimpleTemplate(
	//	"store string",
	//	"AStSt",
	//	1,
	//).
	//	WithInitialParameters(template.Parameters{20}).
	//	WithTransactionEdit(
	//		func(parameters template.Parameters) template.TransactionEdit {
	//			return func(c template.Context, a template.Account) (template.TransactionEdit, error) {
	//				body := fmt.Sprintf(`
	//					signer.storage.save("%s", to: /storage/AStSt)
	//					signer.storage.load<String>(from: /storage/AStSt)
	//				`, stringOfLen(parameters[0]))
	//
	//				return template.TransactionEdit{
	//					PrepareBlock: template.LoopTemplate(100, body),
	//				}, nil
	//			}
	//		},
	//	).
	//	WithAccountSetup(
	//		func(ctx context.Context, c template.Context, interaction template.ChainInteraction, account template.Account, params template.Parameters) error {
	//			// remove anything that might be in storage
	//			return template.RunTransactionBodyAsAccount(
	//				ctx,
	//				"signer.storage.load<String>(from: /storage/AStSt)",
	//				account,
	//				interaction,
	//				nil,
	//			)
	//		},
	//	),
	// no unique intensities
	//template.NewSimpleTemplate(
	//	"create long string",
	//	"LngStr",
	//	1,
	//).
	//	WithInitialParameters(template.Parameters{100}).
	//	WithTransactionEdit(
	//		func(parameters template.Parameters) template.TransactionEdit {
	//			return func(c template.Context, a template.Account) (template.TransactionEdit, error) {
	//				body := fmt.Sprintf(`
	//					var s = "%s"
	//					var S = ""
	//                    %s
	//				`, stringOfLen(parameters[0]+1), template.LoopTemplate(100,
	//					`
	//						S = S.concat(s)
	//					`))
	//
	//				return template.TransactionEdit{
	//					PrepareBlock: body,
	//				}, nil
	//			}
	//		},
	//	),
	// no unique intensities
	//template.NewSimpleTemplate(
	//	"create long string with templating",
	//	"LngStrTmpl",
	//	1,
	//).
	//	WithInitialParameters(template.Parameters{100}).
	//	WithTransactionEdit(
	//		func(parameters template.Parameters) template.TransactionEdit {
	//			return func(c template.Context, a template.Account) (template.TransactionEdit, error) {
	//				body := fmt.Sprintf(`
	//					var s = "%s"
	//					var S = ""
	//                    %s
	//				`, stringOfLen(parameters[0]+1), template.LoopTemplate(100,
	//					`
	//						S = "\(S)\(s)"
	//					`))
	//
	//				return template.TransactionEdit{
	//					PrepareBlock: body,
	//				}, nil
	//			}
	//		},
	//	),
	template.NewSimpleTemplate(
		"borrow string",
		"ABrSt",
		1,
	).
		WithInitialParameters(template.Parameters{1326}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				return template.TransactionEdit{
					PrepareBlock: `
					let strings = signer.storage.borrow<&[String]>(from: /storage/ABrSt)!
					var i = 0
					var lenSum = 0
					while (i < strings.length) {
						lenSum = lenSum + strings[i].length
						i = i + 1
					}
					`,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"copy string",
		"ACpSt",
		1,
	).
		WithInitialParameters(template.Parameters{1352}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				return template.TransactionEdit{
					PrepareBlock: `
					let strings = signer.storage.copy<[String]>(from: /storage/ACpSt)!
					var i = 0
					var lenSum = 0
					while (i < strings.length) {
						lenSum = lenSum + strings[i].length
						i = i + 1
					}
					`,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"copy string and save a duplicate",
		"ACpStSv",
		1,
	).
		WithInitialParameters(template.Parameters{1223}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				return template.TransactionEdit{
					PrepareBlock: `
					let strings = signer.storage.copy<[String]>(from: /storage/ACpStSv)!
					var i = 0
					var lenSum = 0
					while (i < strings.length) {
						lenSum = lenSum + strings[i].length
						i = i + 1
					}
					signer.storage.save(strings, to: /storage/ACpStSv2)
					`,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"store and load dict string",
		"AStDSt",
		1,
	).
		WithInitialParameters(template.Parameters{786}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						signer.storage.save<{String: String}>(%s, to: /storage/AStDSt)
						signer.storage.load<{String: String}>(from: /storage/AStDSt)
					`, stringDictOfLen(parameters[0], 75))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"store load and destroy dict string",
		"ALdDStD",
		1,
	).
		WithInitialParameters(template.Parameters{3324}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := `
						let strings = signer.storage.load<{String: String}>(from: /storage/ALdDStD)!
						for key in strings.keys {
							strings.remove(key: key)
						}
					`

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"borrow dict string",
		"ABrDSt",
		1,
	).
		WithInitialParameters(template.Parameters{206}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				return template.TransactionEdit{
					PrepareBlock: `
					let strings = signer.storage.borrow<&{String: String}>(from: /storage/ABrDSt)!
					var lenSum = 0
					strings.forEachKey(fun (key: String): Bool {
						lenSum = lenSum + strings[key]!.length
						return true
					})
					`,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"copy dict string",
		"ACpDSt",
		1,
	).
		WithInitialParameters(template.Parameters{813}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				return template.TransactionEdit{
					PrepareBlock: `
					let strings = signer.storage.copy<{String: String}>(from: /storage/ACpDSt)!
					var lenSum = 0
					strings.forEachKey(fun (key: String): Bool {
						lenSum = lenSum + strings[key]!.length
						return true
					})
					`,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"copy string dict and save a duplicate",
		"ACpDStSv",
		1,
	).
		WithInitialParameters(template.Parameters{1179}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				return template.TransactionEdit{
					PrepareBlock: `
					let strings = signer.storage.copy<{String: String}>(from: /storage/ACpDStSv)!
					var lenSum = 0
					strings.forEachKey(fun (key: String): Bool {
						lenSum = lenSum + strings[key]!.length
						return true
					})
					signer.storage.save(strings, to: /storage/ACpDStSv2)
					`,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"load dict and destroy it",
		"DestDict",
		1,
	).
		WithInitialParameters(template.Parameters{967}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				return template.TransactionEdit{
					PrepareBlock: `
					let r <- signer.storage.load<@{String: AnyResource}>(from: /storage/DestDict)!
					destroy r
					`,
				}, nil
			},
		),
	simpleTemplateWithLoop(
		"add key to account",
		"KA",
		49,
		`
				let key = PublicKey(
					publicKey: "f7901e9161b9b53f2e1f27b0f1e4711fcc8f234a90f55fd2068a67b152948389c0ee1e40f74a0e194ef7c2b59666270b16d52cf585fd8e65fc00958f78af77b0".decodeHex(),
					signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1
				)
		
				signer.keys.add(
					publicKey: key,
					hashAlgorithm: HashAlgorithm.SHA3_256,
					weight: 0.0
				)
			`,
	),
	simpleTemplateWithLoop(
		"add and revoke key to account",
		"KAR",
		38,
		`
				let key = PublicKey(
					publicKey: "f7901e9161b9b53f2e1f27b0f1e4711fcc8f234a90f55fd2068a67b152948389c0ee1e40f74a0e194ef7c2b59666270b16d52cf585fd8e65fc00958f78af77b0".decodeHex(),
					signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1
				)
		
				let ac = signer.keys.add(
					publicKey: key,
					hashAlgorithm: HashAlgorithm.SHA3_256,
					weight: 0.0
				)
				signer.keys.revoke(keyIndex: ac.keyIndex)
			`,
	),
	simpleTemplateWithLoop(
		"get account key",
		"KGet",
		177,
		`
				let key = signer.keys.get(keyIndex: 0)
			`,
	),
	template.NewSimpleTemplate(
		"Get contracts",
		"GetCon",
		1,
	).
		WithInitialParameters(template.Parameters{65}).
		WithTransactionEdit(
			func(params template.Parameters) (template.TransactionEdit, error) {
				body := `
						signer.contracts.names
					`

				return template.TransactionEdit{
					PrepareBlock: template.LoopTemplate(params[0], body),
				}, nil
			},
		),

	template.NewSimpleTemplate(
		"hash",
		"H",
		1,
	).
		WithInitialParameters(template.Parameters{328}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						Crypto.hash("%s".utf8, algorithm: HashAlgorithm.SHA2_256)
					`, stringOfLen(20))

				return template.TransactionEdit{
					PrepareBlock: template.LoopTemplate(parameters[0], body),
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"String toLower",
		"STL",
		1,
	).
		WithInitialParameters(template.Parameters{5221}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						var s = "%s"
						s = s.toLower()
					`, stringOfLen(parameters[0]))

				return template.TransactionEdit{
					PrepareBlock: template.LoopTemplate(500, body),
				}, nil
			},
		),

	simpleTemplateWithLoop(
		"Get Current Block",
		"GCB",
		1107,
		`
				getCurrentBlock()
			`,
	),
	simpleTemplateWithLoop(
		"Get Block At",
		"GBA",
		548,
		`
				let at = getCurrentBlock().height
				getBlock(at: at)
			`,
	),
	simpleTemplateWithLoop(
		"Destroy Resource Dictionary",
		"DRD",
		1058,
		`
				let r: @{String: AnyResource} <- {}
				destroy r
			`,
	),
	simpleTemplateWithLoop(
		"Parse UFix64",
		"PUFix",
		2032,
		`
				let smol: UFix64? = UFix64.fromString("0.123456")
			`,
	),
	simpleTemplateWithLoop(
		"Parse Fix64",
		"PFix",
		2047,
		`
				let smol: Fix64? = Fix64.fromString("-0.123456")
			`,
	),
	simpleTemplateWithLoop(
		"Parse UInt64",
		"PUInt64",
		2578,
		`
				let smol: UInt64? = UInt64.fromString("123456")
			`,
	),
	simpleTemplateWithLoop(
		"Parse Int64",
		"PInt64",
		2417,
		`
				let smol: Int64? = Int64.fromString("-123456")
			`,
	),
	simpleTemplateWithLoop(
		"Parse Int",
		"PInt",
		2460,
		`
				let smol: Int? = Int.fromString("-12345")
			`,
	),
	simpleTemplateWithLoop(
		"Issue storage capability",
		"ISCap",
		182,
		`
				let cap = signer.capabilities.storage.issue<&Int>(/storage/foo)
			`,
	),
	simpleTemplateWithLoop(
		"Get key count",
		"GKC",
		677,
		`
				let count = signer.keys.count
			`,
	),
	template.NewSimpleTemplate(
		"Create Key ECDSA_P256",
		"CrKeyP256",
		1,
	).
		WithInitialParameters(template.Parameters{112}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				seed := make([]byte, crypto.MinSeedLength)
				for i := range seed {
					seed[i] = 0
				}

				privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
				if err != nil {
					return template.TransactionEdit{}, err
				}
				key := hex.EncodeToString(privateKey.PublicKey().Encode())

				body := fmt.Sprintf(`
						  let publicKey = PublicKey(
							publicKey: "%s".decodeHex(),
							signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
						)
					`, key)

				return template.TransactionEdit{
					PrepareBlock: template.LoopTemplate(parameters[0], body),
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Create Key ECDSA_secp256k1",
		"CrKeysecp256k1",
		1,
	).
		WithInitialParameters(template.Parameters{112}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := `
						let publicKey = PublicKey(
							publicKey: "PUBLIC_KEY_PLACEHOLDER".decodeHex(),
							signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1
						)
					`

				return template.TransactionEdit{
					PrepareBlock: template.LoopTemplate(parameters[0], body),
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Create Key BLS_BLS12_381",
		"CrKeyBLS",
		1,
	).
		WithInitialParameters(template.Parameters{81}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := `
						let publicKey = PublicKey(
							publicKey: "PUBLIC_KEY_PLACEHOLDER".decodeHex(),
							signatureAlgorithm: SignatureAlgorithm.BLS_BLS12_381
						)
					`

				return template.TransactionEdit{
					PrepareBlock: template.LoopTemplate(parameters[0], body),
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Array Insert",
		"ArrIns",
		1,
	).
		WithInitialParameters(template.Parameters{1577}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						let x = [0]
						%s
					`, template.LoopTemplate(
					parameters[0],
					`x.insert(at: i, 1)`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Array Insert Remove",
		"ArrInsDel",
		1,
	).
		WithInitialParameters(template.Parameters{1232}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						let x = [0]
						%s
					`, template.LoopTemplate(
					parameters[0],
					`
							x.insert(at: 0, 1)
							x.remove(at: 1)
							`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Array Insert Set Remove",
		"ArrInsSetDel",
		1,
	).
		WithInitialParameters(template.Parameters{994}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						let x = [0]
						%s
					`, template.LoopTemplate(
					parameters[0],
					`
							x.insert(at: 0, 1)
							x[0] = i
							x.remove(at: 1)
							`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Array Insert Map",
		"ArrInsMap",
		1,
	).
		WithInitialParameters(template.Parameters{1345}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						let x = [0]
						%s
						let addOne =
							fun (_ v: Int): Int {
								return v+1
							}
						let y = x.map(addOne)
					`, template.LoopTemplate(
					parameters[0],
					`
							x.insert(at: 0, i)
							`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Array Insert Filter",
		"ArrInsFilt",
		1,
	).
		WithInitialParameters(template.Parameters{1356}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						let x = [0]
						%s
						let isEven =
							view fun (element: Int): Bool {
								return element %% 2 == 0
							}
						let y = x.filter(isEven)
					`, template.LoopTemplate(
					parameters[0],
					`
							x.insert(at: 0, i)
							`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Array Append",
		"ArrApp",
		1,
	).
		WithInitialParameters(template.Parameters{1905}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
				let x = [0]
				%s
			`, template.LoopTemplate(
					parameters[0],
					`
					x.append(i)
					`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Dict Insert",
		"DictIns",
		1,
	).
		WithInitialParameters(template.Parameters{1598}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						let x = {"0": 0}
						%s
					`, template.LoopTemplate(
					parameters[0],
					`
							x.insert(key: i.toString(), i)
							`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Dict Insert Remove",
		"DictInsDel",
		1,
	).
		WithInitialParameters(template.Parameters{713}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						   let x = {"0": 0}
                           %s
					`, template.LoopTemplate(
					parameters[0],
					`
							x.insert(key: i.toString(), i)
							x.remove(key: (i-1).toString())
							`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Dict Insert Set Remove",
		"DictInsSetDel",
		1,
	).
		WithInitialParameters(template.Parameters{565}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						let x = {"0": 0}
						%s
					`, template.LoopTemplate(
					parameters[0],
					`
							x.insert(key: i.toString(), i)
							x[(i-1).toString()] = i
							x.remove(key: (i-1).toString())
							`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Dict Iter Copy",
		"DictItrCpy",
		1,
	).
		WithInitialParameters(template.Parameters{667}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				body := fmt.Sprintf(`
						let x = {"0": 0}
						let y = {"0": 0}
						%s
						x.forEachKey(fun (key: String): Bool {
							y[key] = x[key]
							return true
						})
					`, template.LoopTemplate(
					parameters[0],
					`
							x.insert(key: i.toString(), i)
							`,
				))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Array Create Batch",
		"ArrCB",
		1,
	).
		WithInitialParameters(template.Parameters{226}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				sumStr := "0"
				for i := 0; i < int(parameters[0]); i++ {
					sumStr += fmt.Sprintf(",%d", i)
				}

				body := fmt.Sprintf(`
						var i = 0
						while i < 200 {
							i = i + 1
							let a = [%s]
						}
					`, sumStr)

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),

	template.NewSimpleTemplate(
		"Verify Signature",
		"VerSig",
		1,
	).
		WithInitialParameters(template.Parameters{14}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				numKeys := parameters[0] + 1
				message := []byte("hello world")

				rawKeys := make([]string, numKeys)
				signers := make([]crypto.Signer, numKeys)

				for i := 0; i < int(numKeys); i++ {
					seed := make([]byte, crypto.MinSeedLength)
					_, err := rand.Read(seed)
					if err != nil {
						return template.TransactionEdit{}, err
					}

					privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
					if err != nil {
						return template.TransactionEdit{}, err
					}
					rawKeys[i] = hex.EncodeToString(privateKey.PublicKey().Encode())
					sig, err := crypto.NewInMemorySigner(privateKey, crypto.SHA3_256)
					if err != nil {
						return template.TransactionEdit{}, err
					}
					signers[i] = sig
				}

				keyListAdd := ""
				for i := 0; i < int(numKeys); i++ {
					keyListAdd += fmt.Sprintf(`
						keyList.add(
							PublicKey(
								publicKey: "%s".decodeHex(),
								signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
							),
							hashAlgorithm: HashAlgorithm.SHA3_256,
							weight: 1.0/%d.0+0.000001 ,
						)
						`, rawKeys[i], int(numKeys),
					)
				}

				signaturesAdd := ""
				for i := 0; i < int(numKeys); i++ {
					signaturesAdd += fmt.Sprintf(`
						signatureSet.append(
							Crypto.KeyListSignature(
								keyIndex: %d,
								signature: "%s".decodeHex()
							)
						)
						`, i, fmt.Sprintf("SIGNATURE_PLACEHOLDER_%d", i),
					)
				}

				body := fmt.Sprintf(`
						let keyList = Crypto.KeyList()

						%s

						let signatureSet: [Crypto.KeyListSignature] = []
		
						%s
		
						let domainSeparationTag = "FLOW-V0.0-user"
						let message = "%s".decodeHex()
						
						
						let valid = keyList.verify(
							signatureSet: signatureSet,
							signedData: message,
							domainSeparationTag: domainSeparationTag
						)
						if !valid {
							panic("invalid signature")
						}
					`, keyListAdd, signaturesAdd, hex.EncodeToString(message))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Aggregate BLS aggregate signature",
		"BLSAggSig",
		1,
	).
		WithInitialParameters(template.Parameters{186}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				numSigs := int(parameters[0]) + 1
				input := make([]byte, 100)
				_, err := rand.Read(input)
				if err != nil {
					return template.TransactionEdit{}, err
				}

				signatures := ""
				for i := 0; i < numSigs; i++ {
					signatures += fmt.Sprintf(`
									signatures.append("%s".decodeHex())
								`, fmt.Sprintf("SIGNATURE_PLACEHOLDER_%d", i))
				}

				body := fmt.Sprintf(`
							var signatures: [[UInt8]] = []
							%s
							BLS.aggregateSignatures(signatures)!
						`, signatures)

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"Aggregate BLS aggregate keys",
		"BLSAggKey",
		1,
	).
		WithInitialParameters(template.Parameters{41}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				numSigs := int(parameters[0]) + 1
				pks := make([]crypto2.PublicKey, 0, numSigs)
				signatureAlgorithm := crypto2.BLSBLS12381
				input := make([]byte, 100)
				_, err := rand.Read(input)
				if err != nil {
					return template.TransactionEdit{}, err
				}

				for i := 0; i < numSigs; i++ {
					seed := make([]byte, crypto2.KeyGenSeedMinLen)
					_, err := rand.Read(seed)
					if err != nil {
						return template.TransactionEdit{}, err
					}
					sk, err := crypto.GeneratePrivateKey(signatureAlgorithm, seed)
					if err != nil {
						return template.TransactionEdit{}, err
					}

					pks = append(pks, sk.PublicKey())
				}

				pkString := ""
				for i := 0; i < numSigs; i++ {
					pkString += fmt.Sprintf(`
									pks.append(PublicKey(
										publicKey: "%s".decodeHex(), 
										signatureAlgorithm: SignatureAlgorithm.BLS_BLS12_381
									))
								`, hex.EncodeToString(pks[i].Encode()))
				}

				body := fmt.Sprintf(`
							let pks: [PublicKey] = []
							%s
							BLS.aggregatePublicKeys(pks)!.publicKey
						`, pkString)

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"BLS verify signature",
		"BLSVer",
		1,
	).
		WithInitialParameters(template.Parameters{32}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				numSigs := int(parameters[0]) + 1
				pks := make([]crypto2.PublicKey, 0, numSigs)
				input := make([]byte, 100)
				_, err := rand.Read(input)
				if err != nil {
					return template.TransactionEdit{}, err
				}

				message := []byte("random_message")

				signaturesString := ""
				for i := 0; i < numSigs; i++ {
					signaturesString += fmt.Sprintf(`
						signatures.append("%s".decodeHex())
					`, fmt.Sprintf("SIGNATURE_PLACEHOLDER_%d", i))
				}

				pkString := ""
				for i := 0; i < numSigs; i++ {
					pkString += fmt.Sprintf(`
									pks.append(PublicKey(
										publicKey: "%s".decodeHex(), 
										signatureAlgorithm: SignatureAlgorithm.BLS_BLS12_381
									))
								`, hex.EncodeToString(pks[i].Encode()))
				}

				body := fmt.Sprintf(`
							var pks: [PublicKey] = []
							var signatures: [[UInt8]] = []
							%s
							%s
							let aggPk = BLS.aggregatePublicKeys(pks)!
							let aggSignature = BLS.aggregateSignatures(signatures)!
							let boo = aggPk.verify(
										signature: aggSignature, 
										signedData: "%s".decodeHex(),
										domainSeparationTag: "random_tag", 
										hashAlgorithm: HashAlgorithm.KMAC128_BLS_BLS12_381)
							if !boo {
								panic("invalid signature")
							}
						`, pkString, signaturesString, hex.EncodeToString(message))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
	template.NewSimpleTemplate(
		"BLS verify proof of possession",
		"BLSVerPoP",
		1,
	).
		WithInitialParameters(template.Parameters{8}).
		WithTransactionEdit(
			func(parameters template.Parameters) (template.TransactionEdit, error) {
				signatureAlgorithm := crypto2.BLSBLS12381
				seed := make([]byte, crypto2.KeyGenSeedMinLen)
				_, err := rand.Read(seed)
				if err != nil {
					return template.TransactionEdit{}, err
				}
				sk, err := crypto.GeneratePrivateKey(signatureAlgorithm, seed)
				if err != nil {
					return template.TransactionEdit{}, err
				}
				pk := sk.PublicKey()

				proof, err := crypto2.BLSGeneratePOP(sk)
				if err != nil {
					return template.TransactionEdit{}, err
				}

				body := fmt.Sprintf(`
							let p = PublicKey(
								publicKey: "%s".decodeHex(), 
								signatureAlgorithm: SignatureAlgorithm.BLS_BLS12_381
							)
							var proof = "%s".decodeHex()

							%s
						`,
					hex.EncodeToString(pk.Encode()),
					hex.EncodeToString(proof.Bytes()),
					template.LoopTemplate(parameters[0], `
							var valid = p.verifyPoP(proof)
							if !valid {
								panic("invalid proof of possession")
							}
						`))

				return template.TransactionEdit{
					PrepareBlock: body,
				}, nil
			},
		),
}

func stringOfLen(length uint64) string {
	someString := make([]byte, length)
	for i := 0; i < len(someString); i++ {
		someString[i] = 'x'
	}
	return string(someString)
}

func stringDictOfLen(dictLen uint64, stringLen uint64) string {
	builder := strings.Builder{}
	builder.WriteRune('{')
	for i := uint64(0); i < dictLen; i++ {
		if i > 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune('"')
		someString := make([]byte, stringLen)
		for i := 0; i < len(someString); i++ {
			someString[i] = 'x'
		}
		builder.WriteString(string(someString))
		builder.WriteString(strconv.Itoa(int(i)))
		builder.WriteRune('"')
		builder.WriteRune(':')
		builder.WriteRune('"')
		builder.WriteString(string(someString))
		builder.WriteRune('"')
	}
	builder.WriteRune('}')
	return builder.String()
}
