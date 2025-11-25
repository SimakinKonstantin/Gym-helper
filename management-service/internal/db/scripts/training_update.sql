UPDATE trainings
SET name=$1, exercises=$2
WHERE id=$3 and user_login=$4;