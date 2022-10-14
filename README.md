# LeiShen

Leishen is a flash loan attack detection tool.

You can enter the hash value of a transaction. Leishen will analyze all its internal transactions and ERC20 transfer records to judge whether the transaction is a flash loan attack.

# Prerequisite

- Go environment

# Compile

> `git clone https://github.com/tony4paper/LeiShen.git`
>
> `make`
>
> `build/bin/LeiShen predicate <TxHash>`

# Output

The results will be output to `RST.csv`. Header is:

`block_number,tx_hash,time,internal_records,from,to,is_attack,pair,max_min,all_swap,rate`

- block_number: Height of the block where the transaction is located
- tx_hash: Hash value of the transaction
- time: Block generation time
- internal_records: Number of internal records
- from: Account address of the transaction sender
- to: Account address of the transaction recipient
- is_attack: The Transaction is a flash loan attack
- pair: Involved token pairs
- max_min: Swap with maximum and minimum exchange rate
- all_swap: All swaps that cause exchange rate fluctuations
- rate: Maximum exchange rate fluctuation

# Result

`flpAttacks/44attacks.xlsx` shows all newly-detected flash loan attacks in Ethereum.
