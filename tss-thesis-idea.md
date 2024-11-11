# Threshold Signature Schemes in Practice: Implementing a Secure Transport Layer for TSS in Golang

## Goals

1. Highlight the current THORChain TSS problems such as
    - Relying on the old vulnerable dependencies. Moreover, most of these dependencies are not backward-compatible, so
      updating would require changes in multiple repositories.
    - Relying on the unaudited TSS library with GG20 (Gennaro and Goldfeder CCS 2020)
    - GG20 is unsupported by most of the common libraries (they are usually based on GG18), thus the dependency is
      unique and irreplaceable
    - Lack of documentation thus increased bus factor and transparency of the solution
    - Still existing problems listed in the 2020 Kudelski Security audit

2. Create a comparison table among TSS cryptography libraries for Golang

3. Prepare and well-document a transport layer on top of TSS cryptography library

4. Add support for key generation and key signing operations

5. Create an open-source library that implements a communication (transport) layer on top of the TSS cryptography using
   a go-libp2p library

6. Benchmark the library and compare against the existing THORChain TSS library

7. Additionally, briefly introduce a TSS technology:
    - Main concepts
    - Field of applications
    - Real-world examples
    - Public key cryptography
    - Multi-party computation
    - Shamir secret sharing and Multisig

## References

| Title                                                                                          | Link                                                                                                      |
|------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------|
| THORChain TSS paper                                                                            | [Link](https://github.com/thorchain/Resources/blob/master/Whitepapers/THORChain-TSS-Paper-June2020.pdf)    |
| THORChain TSS benchmark                                                                        | [Link](https://github.com/thorchain/Resources/blob/master/Whitepapers/THORChain-TSS-Benchmark-July2020.pdf)|
| THORChain TSS transport implementation                                                         | [Link](https://gitlab.com/thorchain/tss/go-tss/-/tree/v1.6.5)                                              |
| PR that moved go-tss to a monorepo                                                             | [Link](https://gitlab.com/thorchain/thornode/-/merge_requests/3696)                                        |
| OKX repo introduction to TSS                                                                   | [Link](https://github.com/okx/threshold-lib/blob/main/docs/Threshold_Signature_Scheme.md)                  |
| R. Gennaro, S. Goldfeder, “Fast Multiparty Threshold ECDSA with Fast Trustless Setup”, 2019    | [Link](https://eprint.iacr.org/2019/114.pdf)                                                               |
| CLI and transport for TSS from Binance                                                         | [Link](https://github.com/bnb-chain/tss/tree/master)                                                       |
| User Guide of Threshold Signature Scheme (TSS) in Binance Chain                                | [Link](https://github.com/bnb-chain/tss/blob/master/doc/UserGuide.md#user-guide-of-threshold-signature-scheme-tss-in-binance-chain) |
| THORChain TSS Kudelski Security audit                                                          | [Link](https://kudelskisecurity.com/wp-content/uploads/ThorchainTSSSecurityAudit.pdf)                      |
| Binance: Threshold Signatures Explained                                                        | [Link](https://academy.binance.com/en/articles/threshold-signatures-explained)                             |
| Binance: What Is a Multisig Wallet?                                                            | [Link](https://academy.binance.com/en/articles/what-is-a-multisig-wallet)                                  |
| Binance: TSS library                                                                           | [Link](https://github.com/bnb-chain/tss-lib)                                                               |

## Library selection (WIP)

| Repo                                                                      | Keygen | Signing | Transport* | Security                | Last release | Go version | Notes                  |
|---------------------------------------------------------------------------|--------|---------|------------|-------------------------|--------------|------------|------------------------|
| [bnb-chain/tss-lib](https://github.com/bnb-chain/tss-lib)                 | ✅      | ✅       | ❌          | Audited on Oct 10, 2019 | Jan 16, 2024 | 1.16       | 705 stars              |
| [thorchain/tss](https://gitlab.com/thorchain/tss/go-tss)                  | ✅      | ✅       | ✅          | Audited on Jun 16, 2020 | Fer 8, 2024  | 1.20       | Production-use example |
| [getamis/alice](https://github.com/getamis/alice)                         | ✅      | ✅       | ❌          | Audited on May 19, 2020 | Nov 30, 2023 | 1.20       | Granted by Coinbase    |
| [taurusgroup/frost-ed25519](https://github.com/taurusgroup/frost-ed25519) | ✅      | ✅       | ❌          | Not audited             | Mar 11, 2021 | 1.14       | Good README            |
| [unit410/threshold-ed25519](https://gitlab.com/unit410/threshold-ed25519) | ✅      | ✅       | ❌          | Not audited             | Feb 21, 2020 | 1.19       |                        |
| [entropyxyz/synedrion](https://github.com/entropyxyz/synedrion)           |        |         |            |                         |              |            | TODO                   |
| [coinbase/kryptology](https://github.com/coinbase/kryptology)             |        |         |            | Papers + HackerOne      | Dec 20, 2021 | 1.17       | Archived               |
| [SwingbyProtocol/tss-lib](https://github.com/SwingbyProtocol/tss-lib)     |        |         |            |                         |              |            | Fork of binance        |

__*__ Transport refers to the parties communication during TSS signing.

### bnb-chain/tss-lib

Pros:

* Was [audited](https://github.com/bnb-chain/tss-lib?tab=readme-ov-file#security-audit) on October 10, 2019, by the
  Kudelski Security
* 700+ stars
* A lot of contributors
* Many libs use it as a basis
* Actively maintained

Cons:

* Doesn't have built-in transport (but there are repos with its implementation from bnb-chain;
  see [References](#references) for details)
* Old Go version

### thorchain/tss

Pros:

* Was [audited](https://kudelskisecurity.com/wp-content/uploads/ThorchainTSSSecurityAudit.pdf) on June 16, 2020, by the
  Kudelski Security
* Has its own transport
* Actively maintained (11 contributors committing periodically)
* [Production-ready example](https://gitlab.com/thorchain/thornode/-/tree/develop/bifrost/tss)
* Has a built-in leader election

Cons:

* Not popular (11 contributors, 6 stars)
* Doubts on the quality of code
* Hard to import (or even impossible without modifications) due to legacy dependencies

### getamis/alice

Pros:

* Was [audited](https://github.com/getamis/alice?tab=readme-ov-file#audit-report) on May 19, 2020, by the Kudelski
  Security
* 340+ stars
* Actively maintained
* Wide range of cryptographic libs (meaning maintainers know what they are doing)
* Granted by Coinbase

Cons:

* Doesn't have transport or leader election
* HTSS differs from TSS, will need additional time to dig into it
