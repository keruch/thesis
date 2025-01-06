
# Threshold Signature Schemes in Practice: Implementing a Secure Transport Layer for TSS in Golang

## Abstract

This paper aims on the implementation of a secure transport layer for Threshold Signature Schemes (TSS) in Golang. There are several implementations of the TSS cryptography in Golang (for example, go-tss lib from Binance[7]), but the main problem is that there are no libraries in Golang that would have the transport layer on top of TSS. Transport layer (P2P) is a crucial part of the TSS[8], as it is responsible for the communication between the parties during the signing process. The only library that have the transport layer is THORChain TSS, which is built on top of the modified Binance library for TSS cryptography. However, THORChain TSS it is not well-documented, maintained by a small team, hosted in GitLab (which is not convenient for users in general, and this library is hard to find with random walks), uses outdated and not backward-compatible dependencies (so it's not compatible with applications using modern dependencies), and has some security and reliability concerns. Moreover, it has been archived as a library and moved to THORChain monorepo[6], making it entirely unmaintainable as the external library. Despite the above, lack of alternatives forces developers to use THORChain TSS (e.g., ZetaChain[1], Maya Protocol[2], SwipeWallet[3]), or create their own solutions (e.g., Sygma (ex. ChainSafe)[4], Threshold Network[5]).

This paper introduces an implementation of the secure communication protocol for TSS signing in Golang and benchmarks it against the existing THORChain solution. Importantly, this paper leverages the existing TSS cryptography libraries, building the communication layer on top of them. It does not roll its own cryptography but rather relies on existing audited solutions. On top of that, this paper provides a comparison of the existing TSS cryptography libraries in Golang.

## Objectives

1. Introduce TSS Technology and Applications
    - Introduce glossary terms.
    - Provide an overview of the TSS technology, its main concepts, and applications.
    - Present real-world examples of TSS usage.
    - Explain the importance of TSS in public key cryptography and multi-party computation.

2. Address Limitations of Existing Solutions
    - Indicate drawbacks of the THORChain TSS implementation.
    - List alternatives that different projects use.
    - Compare existing TSS cryptography libraries in Golang.

3. Design and Implement a Secure Transport Layer for TSS in Golang
    - Develop a communication protocol specifically tailored for TSS operations: key generation, key signing.
    - Select the most suitable TSS cryptography library.

4. Benchmark and Compare Implementations
    - Conduct performance evaluations comparing with the existing THORChain solution.
    - Analyze time and memory utilization for key generation and key signing depending on the different number of parties.

5. Enhance Security and Reliability
    - Employ techniques from Binance Guide[8] to ensure secure communication.
    - (???) Ensure the transport layer is resilient against common network attacks and failures.

6. Contribute to the Community
    - Release the implementation as an open-source project.
    - Provide comprehensive documentation and examples.
    - (???) Try to make the library compatible with the existing projects.
    - (???) Create a demo application that uses the library.

7. Introduce Further Improvements
    - Everything marked with (???) in the previous objectives.

## Execution Plan

### Initial Research and Planning (Weeks 1-2, November)

1. Literature Review:
    - Study foundational papers and resources on TSS, multi-party computation, and related cryptographic concepts.
    - Compile a glossary of key terms to be used throughout the thesis.
    - Prepare and analyze a list of references.
2. Documentation Drafting:
    - Write the introductory sections, covering:
        - Overview of TSS technology and its importance in public key cryptography.
        - Real-world applications and examples of TSS usage.
        - Explanation of how TSS enhances security in distributed systems.

### Analysis of Existing Solutions (Weeks 3-4, December)

1. Examine THORChain TSS Implementation:
    - Review the codebase and documentation (or lack thereof).
    - Identify specific issues related to maintenance, outdated dependencies, and security concerns.
2. Research Alternatives:
    - Investigate other projects (e.g., Sygma, Threshold Network) and their approaches to TSS.
    - Compile a list of existing TSS cryptography libraries in Golang.
3. Comparative Analysis:
    - Evaluate the features, pros, and cons of each library.
    - Prepare a comparison chart or table for clarity.

### Design of the Secure Transport Layer (Weeks 5-6, December)

1. Select TSS Cryptography Library:
    - Choose the most suitable existing TSS library (e.g., Binance's go-tss) to build upon.
2. Protocol Design:
    - Define the requirements for the communication protocol specific to TSS operations like key generation and key signing.
    - Design message formats, communication sequences, and error-handling mechanisms.
    - Incorporate security measures based on best practices and Binance Guide[8].
3. Architectural Planning:
    - Outline the architecture of the transport layer, ensuring scalability (it may have multiple parties).
    - Plan how the transport layer will interface with the chosen TSS cryptography library.

### Implementation of the Transport Layer (Weeks 7-10, January)

1. Coding:
    - Implement the designed communication protocol.
    - Ensure secure communication channels, possibly using TLS or other encryption methods.
    - Develop modules for key generation and key signing processes.
    - Create APIs (HTTP or gRPC server and simple CLI) to allow users run the environment manually.
2. Testing:
    - Write unit tests for individual components.

### Testing and Validation (Weeks 11-12, February)

1. Security Enhancements:
    - Apply techniques from the Binance Guide [8] to strengthen security.
2. Resilience Measures (Optional):
    - (If time permits) Ensure the transport layer can handle network attacks and failures gracefully.
    - Implement retry mechanisms and fault tolerance features.
3. Validation:
    - Perform integration tests to ensure the library works with multiple nodes.

### Benchmarking and Comparison (Weeks 13-14, February)

1. Set Up Benchmarking Environment:
    - Prepare testing scenarios with varying numbers of parties.
2. Performance Evaluation:
    - Measure time and memory utilization during key generation and key signing.
3. Comparison with THORChain TSS:
    - (???) Repeat the same tests using the THORChain TSS implementation.
    - Use THORChain TSS benchmark paper[9] as a reference.
    - Compare results to evaluate improvements or regressions.
4. Data Analysis:
    - Use statistical tools to interpret the benchmarking results.
    - Visualize data using graphs and charts for clarity.

### Documentation and Open-Source Release (Weeks 15-16, March)

1. Documentation:
    - Write comprehensive user guides, including installation instructions and usage examples.
    - Generate API documentation with comments and annotations in the code.
    - Include a clear README.md with project overview and contribution guidelines.
2. Open-Source (Optional, if time permits):
    - Share the project in relevant forums or communities to gather initial feedback.
    - Adjust the library to be compatible with existing projects.
    - Develop a simple demo app showcasing the library's capabilities.

### Thesis Writing and Finalization (Weeks 17-19, April)

1. Compile Thesis Document:
    - Integrate all written sections into a cohesive thesis.
    - Ensure that the introduction, methodology, results, and conclusions are well-articulated.
2. Editing and Proofreading:
    - Review for grammatical errors and clarity.
    - Verify that all citations and references are correctly formatted.
3. Advisor Review:
    - Make necessary revisions based on comments.
4. Preparation for Defense:
    - Prepare presentation slides.
    - Practice the defense speech, anticipate possible questions.

## References

1. [ZetaChain](https://github.com/zeta-chain), [fork of Binance TSS](https://github.com/zeta-chain/tss-lib/tree/threshold-dep-updates)
2. [Maya Protocol](https://gitlab.com/mayachain)
3. [SwipeWallet](https://github.com/SwipeWallet)
4. [Sygma](https://github.com/sygmaprotocol), [relayer code that uses TSS](https://github.com/sygmaprotocol/sygma-relayer)
5. [Threshold Network](https://github.com/threshold-network), [fork of Binance TSS](https://github.com/threshold-network/tss-lib)
6. Issue to move THORChain go-tss to monorepo https://gitlab.com/thorchain/thornode/-/issues/2024
7. Binance TSS https://github.com/bnb-chain/tss-lib
8. How to use TSS securely https://github.com/bnb-chain/tss-lib?tab=readme-ov-file#how-to-use-this-securely
9. THORChain TSS benchmark paper: https://github.com/thorchain/Resources/blob/master/Whitepapers/THORChain-TSS-Benchmark-July2020.pdf

## Additional References

ChatGPT chat with related info: https://chatgpt.com/share/6735ec7c-d8e8-800a-b9a0-0bb3e5f322c5.

| Title                                                                                       | Link                                                                                                                                |
|---------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| THORChain TSS paper                                                                         | [Link](https://github.com/thorchain/Resources/blob/master/Whitepapers/THORChain-TSS-Paper-June2020.pdf)                             |
| THORChain TSS benchmark                                                                     | [Link](https://github.com/thorchain/Resources/blob/master/Whitepapers/THORChain-TSS-Benchmark-July2020.pdf)                         |
| THORChain TSS transport implementation                                                      | [Link](https://gitlab.com/thorchain/tss/go-tss/-/tree/v1.6.5)                                                                       |
| PR that moved go-tss to a monorepo                                                          | [Link](https://gitlab.com/thorchain/thornode/-/merge_requests/3696)                                                                 |
| OKX repo introduction to TSS                                                                | [Link](https://github.com/okx/threshold-lib/blob/main/docs/Threshold_Signature_Scheme.md)                                           |
| R. Gennaro, S. Goldfeder, “Fast Multiparty Threshold ECDSA with Fast Trustless Setup”, 2019 | [Link](https://eprint.iacr.org/2019/114.pdf)                                                                                        |
| CLI and transport for TSS from Binance                                                      | [Link](https://github.com/bnb-chain/tss/tree/master)                                                                                |
| User Guide of Threshold Signature Scheme (TSS) in Binance Chain                             | [Link](https://github.com/bnb-chain/tss/blob/master/doc/UserGuide.md#user-guide-of-threshold-signature-scheme-tss-in-binance-chain) |
| THORChain TSS Kudelski Security audit                                                       | [Link](https://kudelskisecurity.com/wp-content/uploads/ThorchainTSSSecurityAudit.pdf)                                               |
| Binance: Threshold Signatures Explained                                                     | [Link](https://academy.binance.com/en/articles/threshold-signatures-explained)                                                      |
| Binance: What Is a Multisig Wallet?                                                         | [Link](https://academy.binance.com/en/articles/what-is-a-multisig-wallet)                                                           |
| Binance: TSS library                                                                        | [Link](https://github.com/bnb-chain/tss-lib)                                                                                        |
| CHAINBRIDGE MULTI-PARTY SIGNING - RESEARCH                                                  | https://hackmd.io/@timofey/BkG-pTBCK                                                                                                |
| CHAINBRIDGE p2p communication for TSS systems pre-design research                           | https://github.com/ChainSafe/chainbridge-core/issues/279                                                                            |
| CHAINBRIDGE TSS for bridging pre-desing research                                            | https://github.com/ChainSafe/chainbridge-core/issues/280                                                                            |
| Qredo Protocol Yellow Paper                                                                 | https://www.qredo.com/qredo-yellow-paper.pdf                                                                                        |

### Projects Using THORChain TSS

1. https://github.com/SwitchlyProtocol
2. https://github.com/polymerdao
3. https://github.com/Zecrey-Labs
4. https://github.com/0xPellNetwork
5. https://github.com/SwipeWallet
6. https://github.com/zeta-chain
7. https://gitlab.com/mayachain
8. https://int3face.zone/

### Projects Using non-THORChain TSS

1. https://github.com/threshold-network/tss-lib
2. https://github.com/sygmaprotocol

## Remarks

### THORChain TSS Problems

- Relying on the old vulnerable dependencies. Moreover, most of these dependencies are not backward-compatible, so updating would require changes in multiple repositories.
- Relying on the unaudited TSS library with GG20 (Gennaro and Goldfeder CCS 2020)
- GG20 is unsupported by most of the common libraries (they are usually based on GG18), thus the dependency is unique and irreplaceable
- Lack of documentation thus increased bus factor and transparency of the solution
- Still existing problems listed in the 2020 Kudelski Security audit

### TSS Technology Introduction:

- Main concepts
- Field of applications
- Real-world examples
- Public key cryptography
- Multi-party computation
- Shamir secret sharing and Multisig

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

# P2P

Discovering and connecting with other peers is a key challenge in P2P networking. In the past, each P2P application had to develop its own solution for this problem, leading to a lack of reusable, well-documented P2P protocols. IPFS looked to existing research and networking applications for inspiration, but found few code implementations that were usable and adaptable. Many of the existing implementations had poor documentation, restrictive licensing, outdated code, no point of contact, were closed source, deprecated, lacked specifications, had unfriendly APIs, or were tightly coupled with specific use cases and not upgradeable. As a result, developers often had to reinvent the wheel each time they needed P2P protocols, rather than being able to reuse existing solutions.

Security Considerations: https://docs.libp2p.io/concepts/security/security-considerations/
Secure Channels: https://docs.libp2p.io/concepts/secure-comm/overview/
Publish/Subscribe: https://docs.libp2p.io/concepts/pubsub/overview/
P2P launchpad: https://pl-launchpad.io/curriculum/libp2p/objectives/
NAT traversal