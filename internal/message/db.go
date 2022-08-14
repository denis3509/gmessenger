package message

const (
	contactListSQL string = `select u2        as user_id,
 
	"user".username,

	maxTime as max_time,
	(select text_content
	 from message
	 where created_at = maxTime
	   and (sender_id = {:user} or recipient_id = {:user})
	 limit 1) as last_message_text,
	(select is_read  
	 from message
	 where created_at = maxTime
	   and (sender_id = {:user} or recipient_id = {:user})
	 limit 1) as last_message_read,
	unread
from (select (case
			   when (sender_id = {:user})
				   then (sender_id)
			   when (sender_id != {:user})
				   then recipient_id
 end)               as u1,
		  (case
			   when (sender_id = {:user})
				   then (recipient_id)
			   when (sender_id != {:user})
				   then sender_id
			  end)  as u2,

		  max(created_at) as maxTime

   from message
   where (sender_id = {:user} or recipient_id = {:user})
   group by u1, u2
  ) as receivers_id
	  left join "user"
				on receivers_id.u2 = "user".id
	  left join (select sender_id, count(id) as unread
				 from message m
				 where recipient_id = {:user}
				   and is_read = false
				 group by sender_id
) as counter
				on
					u2 = counter.sender_id
order by maxTime desc;`
)
