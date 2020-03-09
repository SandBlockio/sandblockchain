package app

const (
	// Bech32MainPrefix The base prefix for all keys
	Bech32MainPrefix = "sand"

	// PrefixValidator The prefix to add at validators keys
	PrefixValidator = "val"

	// PrefixConensus the prefix to add at consensus keys
	PrefixConsensus = "cons"

	// PrefixPublic the prefix to add at public keys
	PrefixPublic = "pub"

	// PrefixOperator the prefix to add at operator keys
	PrefixOperator = "oper"

	// Bech32PrefixAccAddr defines the prefix of an account's address
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub defines the prefix of an account's public key
	Bech32PrefixAccPub = Bech32MainPrefix + PrefixPublic
	// Bech32PrefixValAddr defines the prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32MainPrefix + PrefixValidator + PrefixOperator
	// Bech32PrefixValPub defines the prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// Bech32PrefixConsAddr defines the prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
	// Bech32PrefixConsPub defines the prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic

)