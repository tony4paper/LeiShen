# LeiShen

Leishen is a detection tool for flash loan based price manipulation attacks (flpAttacks).

You can enter the hash value of a transaction. Leishen will analyze all its internal transactions and ERC20 transfer records to judge whether the transaction is a flpAttack..

# Prerequisite

- Go environment

# Compile

> `git clone https://github.com/tony4paper/LeiShen.git`
>
> `cd LeiShen`
>
> `make`

# Output

The results will be output to `RST.csv`. Header is:

`tx_hash, block_number, from, to, time, types, itx_number, erc_number, elapsed, KRP, MBS, SBS`

- tx_hash: Hash value of the transaction
- block_number: Height of the block where the transaction is located
- from: Account address of the transaction sender
- to: Account address of the transaction recipient
- time: Block generation time
- types: Flash Loan type
- itx_number: Number of internal transactions
- erc_number: Number of ERC20 transfer records
- elapsed: Time elapsed from the transaction to the last transaction
- KRP: Whether the transaction conforms to the attack pattern of Keep Raising Price
- MBS: Whether the transaction conforms to the attack pattern of Multi-Round Buying and Selling
- SBS: Whether the transaction conforms to the attack pattern of Symmetrical Buying and Selling

# Example

`./bin/leishen predicate -b ./dbs/block/ -c ./dbs/contract/ -i ./dbs/internal-transaction/ -r ./dbs/receipt/ -p ./dbs/platform-name/ -fl ./dbs/flash-loan/ -t 0x013be97768b702fe8eccef1a40544d5ecb3c1961ad5f87fee4d16fdc08c78106`

`0x013be97768b702fe8eccef1a40544d5ecb3c1961ad5f87fee4d16fdc08c78106, 10355807, 0xBF675C80540111A310B06e1482f9127eF4E7469A, 0x81D73c55458f024CDC82BbF27468A2dEAA631407, 2020-06-29 02:03:11, dYdX, 2, 315, 75, true, false, true`

# Result

`flpAttacks/44attacks.xlsx` shows our collected 44 flash loan based attacks in empirical study.
