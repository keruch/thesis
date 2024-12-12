# Initial Research and Planning

1. Literature Review:
    - Study foundational papers and resources on TSS, multi-party computation, and related cryptographic concepts.
    - Compile a glossary of key terms to be used throughout the thesis.
    - Prepare and analyze a list of references.
2. Documentation Drafting:
    - Write the introductory sections, covering:
        - Overview of TSS technology and its importance in public key cryptography.
        - Real-world applications and examples of TSS usage.
        - Explanation of how TSS enhances security in distributed systems.

## Glossary

- [ ] TSS 
- [ ] DKG
- [ ] Multi-Party Computation
- [ ] Threshold Cryptography
- [ ] Byzantine Fault Tolerance
- [ ] Secure Multi-Party Computation
- [ ] Verifiable Secret Sharing
- [ ] Homomorphic Encryption
- [ ] Zero-Knowledge Proofs
- [ ] MPC Protocols
- [ ] Distributed Key Generation
- [ ] Key Generation Center
- [ ] Key Share
- [ ] Key Reconstruction
- [ ] Adversary Model
- [ ] Consensus Algorithm
- [ ] Secure Communication Channel
- [ ] Cryptographic Primitives
- [ ] secp256k1: is a specific elliptic curve used in cryptography, most notably in blockchain technologies like Bitcoin and Ethereum. It's not a signature algorithm.
- [ ] ECDSA: Elliptic Curve Digital Signature Algorithm, is a cryptographic algorithm used for creating digital signatures. ECDSA itself does not define which specific curve must be used. It can operate over many elliptic curves, provided they meet certain security and mathematical requirements. One of the mose popular is secp256k1. The signature consists of two numbers, typically denoted as R and S.
- [ ] Ed25519: is a high-performance, secure, and modern public-key signature scheme that is part of the **EdDSA** (Edwards-curve Digital Signature Algorithm) family. Competing with ECDSA, EdDSA is designed to be faster and more secure. The name Ed25519 refers to the specific instantiation of EdDSA using the elliptic curve Curve25519. This is essentially a “locked-in” combination of EdDSA with a chosen curve (Curve25519 in Edwards form).
- [ ] Shamir secret sharing: is a form of secret sharing, where a secret is divided into parts, giving each participant its own unique part. The secret can only be reconstructed when the parts are combined together; individual parts are of no use on their own.
- [ ] Schnorr signature: is a digital signature algorithm that is based on the Schnorr identification protocol. It was the first digital signature scheme to provide a proof of security. It is provably secure in the random oracle model assuming the discrete logarithm problem is hard.

## References

- https://github.com/tnunamak/multisig-decrypt-demo/blob/main/main.go