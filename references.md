# References

Tax:

- "A Computational Approach to Optimal Taxation" by Felix Bierbrauer (2016)
- "Optimal Tax Administration" by Michael Keen and Joel Slemrod (2017)

Crypto:

- "Optimistic Rollups: The Present and Future of Ethereum Scaling" by John Adler, Mikerah Quintyne-Collins, and Alexey
  Akhunov (2020)


- ?? - "On Settlement Finality" by Joachim Zahnentferner (2018)


- Mustafa Al-Bassam, Alberto Sonnino, and Vitalik Buterin. **Fraud and Data Availability Proofs: Maximising Light Client
  Security and Scaling Blockchains with Dishonest Majorities**.
    2018. [Reference 1](https://discovery.ucl.ac.uk/id/eprint/10117245/1/thesis.pdf)
          and [Reference 2](https://arxiv.org/pdf/1809.09044).


- Georgios Konstantopoulos. **(Almost) Everything you need to know about Optimistic Rollup**. In the paradigm.xyz
  website.
    2021. [Reference](https://www.paradigm.xyz/2021/01/almost-everything-you-need-to-know-about-optimistic-rollup).


- Shashank Motepalli, Luciano Freitas, Benjamin Livshits. **SoK: Decentralized Sequencers for Rollups**.
    2023. [Reference](https://arxiv.org/pdf/2310.03616).


- Vitalik Buterin. **An Incomplete Guide to Rollups**. 2021. [Video](https://www.youtube.com/watch?v=wcCHlqgGSH4)
  and [Article](https://vitalik.eth.limo/general/2021/01/05/rollup.html).

  *Simple and understandable overview. Good comparison table between optimistic and ZK rollups.*


- Vitalik Buterin. **The Data Availability Problem.** In Silicon Valley Ethereum Meetup.
    2017. [Video](https://www.youtube.com/watch?v=OJT_fR7wexw).


- Harry Kalodner, Steven Goldfeder, Xiaoqi Chen, S Matthew Weinberg, and Edward W Felten. **Arbitrum: Scalable, private
  smart contracts.** In 27th USENIX Security Symposium (USENIX Security 18), pages 1353–1370,
    2018. [Reference](https://www.usenix.org/system/files/conference/usenixsecurity18/sec18-kalodner.pdf).

  *Interesting thing! Scalability via incentivized verifiers.*


- Lee Bousfield, Rachel Bousfield, Chris Buckland, Ben Burgess, Joshua Colvin, Edward W Felten, Steven Goldfeder, Daniel
  Goldman, Braden Huddleston, Harry Kalodner, et al. **Arbitrum nitro: A secondgeneration optimistic rollup**.
    2022. [Reference](https://github.com/OffchainLabs/nitro/blob/master/docs/Nitro-whitepaper.pdf).

  *Fresh view on the problem.*


- Zhe Ye , Ujval Misra , Jiajun Cheng, Wenyang Zhou, Dawn Song. **Specular: Towards Secure, Trust-minimized Optimistic
  Blockchain Execution**. 2024. [Reference](https://arxiv.org/pdf/2212.05219)

  *Good article with many technical details and comprehensive reference list*


- Daji Landis. **Incentive Non-Compatibility of Optimistic Rollups**.
    2024. [Reference](https://arxiv.org/pdf/2312.01549).

  *Very good paper describing the similar to as mine. Has good references 9, 10, 11.*

Math:

- "A Game-Theoretic Analysis of Inspection Games with Voluntary Disclosure" by Takashi Matsumura and Makoto Shimizu (
    2005)
- "Inspection Games" by Rudolf Avenhaus, Bernhard von Stengel, and Shmuel Zamir (2002)
- "Optimal Verification Procedures" by Yuliy Sannikov (2008)
- "Optimal Auditing with Heterogeneous Audit Perceptions" by Maciej H. Kotowski, David A. Weisbach, and Richard J.
  Zeckhauser (2014)

Glossary:

- Principal-Agent Problem: The relationship between the Layer 1 (L1) blockchain and the Layer 2 (L2) rollup can be
  modeled as a principal-agent problem.
- Equilibrium: In game theory, an equilibrium is a strategy profile where no player has an incentive to deviate.
  Meaning, no player can improve their payoff by changing their strategy while the other players keep their strategies
  fixed.
- Nash Equilibrium: The security of optimistic rollups relies on finding an equilibrium where honest behavior is the
  dominant strategy.

## Reading list

| Priority | Status | Title                                             | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                           | Author, Date          | References                                                                                                              |
|----------|--------|---------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|-------------------------------------------------------------------------------------------------------------------------|
| 1        | ✅      | Incentive Non-Compatibility of Optimistic Rollups | Gives overview of the verifier's dilemma in optimistic rollups. The general result is that optimistic rollups are not incentive-compatible. The security comes with the pure equilibrium, but there are edge cases in the simulation when the undesired equilibrium is reached. I.e., there are values when the sequencer can act dishonestly, while still reaching equilibrium (the outcome is indifferent of the party's strategy while the opponent follows the same strategy) | Daji Landis, 2024     | [Article](https://arxiv.org/pdf/2312.01549)                                                                             |
| 2        | 🌀     | The Data Availability Problem                     |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |                       |                                                                                                                         |
| 3        | ✅      | An Incomplete Guide to Rollups                    | General info about state channels, plasma, optimistic and zk rollups                                                                                                                                                                                                                                                                                                                                                                                                              | Vitalik Buterin, 2021 | [Video](https://www.youtube.com/watch?v=wcCHlqgGSH4) [Article](https://vitalik.eth.limo/general/2021/01/05/rollup.html) |

## TODO

- ill-post problem