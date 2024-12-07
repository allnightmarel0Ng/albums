import streamlit as st

from globals import notification_list

def show_notifications():
    if not "jwt_token" in st.session_state:
        st.error("Notifications are not available without logging in")
        return
    
    for notification in notification_list:
        if notification['success'] == True:
            st.success(notification['message'])
        else:
            st.error(notification['message'])

show_notifications()