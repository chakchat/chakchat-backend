GET /chat/my/all/ 
1. Add chat preview. 
2. new updates aggregated info 
3. think about how chat photo should be returned.
GET /chat/personal/{chatId}
1. Add chat preview
2. new updates aggregated info
3. think about how chat photo should be returned.
PUT /chat/personal/{chatId}/block
1. Should it really be unified with secret personal chat?

Server can response with only last_update_id. 
Client will calculate the range of updates itself.

I think all synchronous calls should be REST calls. 
WebSocket should be used for real-time updates from server.