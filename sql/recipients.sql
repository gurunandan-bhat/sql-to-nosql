select 
	concat('#SET:', s.vRecipientSetName),
	concat('#A:', r.vEmail),
	s.vRecipientSetName,
	r.vEmail,
	r.vFullname,
	r.cStatus
from 
	recipient_set s join
	recipient r on s.iRecipientSetID = r.iRecipientSetID
order by 1, 4
