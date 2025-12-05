UPDATE statistics
SET status = $1
WHERE id=$2 and user_name=$3;
