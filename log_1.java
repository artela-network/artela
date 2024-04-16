5:23PM INF Timed out dur=2950.536547 height=157 module=consensus round=0 step=1
5:23PM DBG received tock height=157 module=consensus round=0 step=1 timeout=2950.536547
5:23PM DBG entering new round current={} height=157 module=consensus round=0
5:23PM DBG entering propose step current={} height=157 module=consensus round=0
5:23PM DBG node is a validator height=157 module=consensus round=0
5:23PM DBG propose step; our turn to propose height=157 module=consensus proposer=3042CCEF91FA9CA42A9CD2B197682832FB4BA84D round=0
-----ReapMaxBytesMaxGas begin
-----all txs: 0
-----all txs is saved from mempool.
-----ReapMaxBytesMaxGas end, len(txs): 0
5:23PM DBG Received tick module=consensus new_ti={"duration":3000000000,"height":157,"round":0,"step":3} old_ti={"duration":2950536547,"height":157,"round":0,"step":1}
5:23PM DBG Timer already stopped module=consensus
5:23PM DBG Scheduled timeout dur=3000 height=157 module=consensus round=0 step=3
5:23PM DBG signed proposal height=157 module=consensus proposal={"Type":32,"block_id":{"hash":"0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057","parts":{"hash":"3EC65D86A42203539144CF557A334C25F18E8559F57813566804A04C72F2F5E6","total":1}},"height":157,"pol_round":-1,"round":0,"signature":"6RdYnsKcqvxcKgkZur5hKKauhayXupGF7ybnjOv8TAVmIErzvmUPZPWOVR5vl/DALAyWKb1Kg17Ln6MlRjpNBA==","timestamp":"2024-01-24T09:23:08.098181286Z"} round=0
5:23PM DBG Broadcast channel=32 module=p2p
5:23PM INF received proposal module=consensus proposal={"Type":32,"block_id":{"hash":"0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057","parts":{"hash":"3EC65D86A42203539144CF557A334C25F18E8559F57813566804A04C72F2F5E6","total":1}},"height":157,"pol_round":-1,"round":0,"signature":"6RdYnsKcqvxcKgkZur5hKKauhayXupGF7ybnjOv8TAVmIErzvmUPZPWOVR5vl/DALAyWKb1Kg17Ln6MlRjpNBA==","timestamp":"2024-01-24T09:23:08.098181286Z"}
{"header":{"version":{"block":11},"chain_id":"artela_11820-1","height":157,"time":"2024-01-24T09:23:03.055502909Z","last_block_id":{"hash":"9D6E2594621D4E9A52A5F12463C092F2307AF5A0FDF6F493F979E75871AE4A4F","parts":{"total":1,"hash":"207DA7197E896AE8B66971894F3564CC4E9A8CC4B38156848C79A02AF3E4BDE0"}},"last_commit_hash":"6E60B294AE391B40B651FE58DCE87BF036B6806A9DD756EBC8878A3B238B3554","data_hash":"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855","validators_hash":"C94F9F7751AB65AF3F5BE83A4CCABD315A74E3EFD5451F594B3CC505E3E2620C","next_validators_hash":"C94F9F7751AB65AF3F5BE83A4CCABD315A74E3EFD5451F594B3CC505E3E2620C","consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F","app_hash":"A45E37493007A1BC06FF53E589B1AE1D1FD459314EDC9E24BF75CC7329AE1872","last_results_hash":"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855","evidence_hash":"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855","proposer_address":"3042CCEF91FA9CA42A9CD2B197682832FB4BA84D"},"data":{"txs":[]},"evidence":{"evidence":[]},"last_commit":{"height":156,"round":0,"block_id":{"hash":"9D6E2594621D4E9A52A5F12463C092F2307AF5A0FDF6F493F979E75871AE4A4F","parts":{"total":1,"hash":"207DA7197E896AE8B66971894F3564CC4E9A8CC4B38156848C79A02AF3E4BDE0"}},"signatures":[{"block_id_flag":2,"validator_address":"3042CCEF91FA9CA42A9CD2B197682832FB4BA84D","timestamp":"2024-01-24T09:23:03.055502909Z","signature":"TgzFot1H6986ioGLVFLJLVn/1hNCW4tfSJsFDIz74XLUocaKJrKI6gpZXwZMjf8ZfOCIyKNcJteA4dv3iaa1Dw=="}]}}
5:23PM INF received complete proposal block hash=0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057 height=157 module=consensus
5:23PM DBG entering prevote step current={} height=157 module=consensus round=0
5:23PM DBG prevote step: ProposalBlock is valid height=157 module=consensus round=0
{"type":1,"height":157,"round":0,"block_id":{"hash":"0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057","parts":{"total":1,"hash":"3EC65D86A42203539144CF557A334C25F18E8559F57813566804A04C72F2F5E6"}},"timestamp":"2024-01-24T09:23:08.12138766Z","validator_address":"3042CCEF91FA9CA42A9CD2B197682832FB4BA84D","validator_index":0,"signature":"Pg8bUkkX1NE5FQZPIattzsMpH8vsr4sW5jGwwPI5QfYkYJx8aWKVf5sbD9qn64Dhaat9ylRVqQjPvkNYhRCYBw=="}
5:23PM DBG signed and pushed vote height=157 module=consensus round=0 vote={"block_id":{"hash":"0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057","parts":{"hash":"3EC65D86A42203539144CF557A334C25F18E8559F57813566804A04C72F2F5E6","total":1}},"height":157,"round":0,"signature":"Pg8bUkkX1NE5FQZPIattzsMpH8vsr4sW5jGwwPI5QfYkYJx8aWKVf5sbD9qn64Dhaat9ylRVqQjPvkNYhRCYBw==","timestamp":"2024-01-24T09:23:08.12138766Z","type":1,"validator_address":"3042CCEF91FA9CA42A9CD2B197682832FB4BA84D","validator_index":0}
5:23PM DBG Broadcast channel=32 module=p2p
5:23PM DBG Attempt to update stats for non-existent peer module=consensus peer=
5:23PM DBG adding vote cs_height=157 module=consensus val_index=0 vote_height=157 vote_type=1
5:23PM DBG Broadcast channel=32 module=p2p
5:23PM DBG added vote to prevote module=consensus prevotes="VoteSet{H:157 R:0 T:SIGNED_MSG_TYPE_PREVOTE +2/3:0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057:1:3EC65D86A422(1) BA{1:x} map[]}" vote={"block_id":{"hash":"0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057","parts":{"hash":"3EC65D86A42203539144CF557A334C25F18E8559F57813566804A04C72F2F5E6","total":1}},"height":157,"round":0,"signature":"Pg8bUkkX1NE5FQZPIattzsMpH8vsr4sW5jGwwPI5QfYkYJx8aWKVf5sbD9qn64Dhaat9ylRVqQjPvkNYhRCYBw==","timestamp":"2024-01-24T09:23:08.12138766Z","type":1,"validator_address":"3042CCEF91FA9CA42A9CD2B197682832FB4BA84D","validator_index":0}
5:23PM DBG updating valid block because of POL module=consensus pol_round=0 valid_round=-1
5:23PM DBG Broadcast channel=32 module=p2p
5:23PM DBG entering precommit step current={} height=157 module=consensus round=0
5:23PM DBG precommit step; +2/3 prevoted proposal block; locking hash=0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057 height=157 module=consensus round=0
{"type":2,"height":157,"round":0,"block_id":{"hash":"0C63241D9C1AE4FB68EC3A8A50A1A5A2D5648BAF80EF22CA4A31C16B47381057","parts":{"total":1,"hash":"3EC65D86A42203539144CF557A334C25F18E8559F57813566804A04C72F2F5E6"}},"timestamp":"2024-01-24T09:23:10.146902609Z","validator_address":"3042CCEF91FA9CA42A9CD2B197682832FB4BA84D","validator_index":0,"signature":"goYNKbiTXmDUt2RptkYIwqeBoHEKWaahsTfcONNT5U92B0HQL4BeyPKqn8YhvkB9mShCVpQOyhg8I/cTcV9+BA=="}

5:23PM INF Timed out dur=3000 height=157 module=consensus round=0 step=3
