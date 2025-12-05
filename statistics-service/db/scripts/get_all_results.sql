SELECT st.id, st.training_id, st.user_name, st.start_time, st.finish_time, st.result_values, st.kcal, st.comment, statuses.name AS status FROM statistics AS st
JOIN statuses ON statuses.id = st.status
WHERE st.user_name=$1;