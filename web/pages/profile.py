import re
import streamlit as st
import requests

from main import chunked, get_image_url, fetch_orders, GATEWAY_PORT

PROFILE_URL = "http://localhost:{GATEWAY_PORT}/profile"
LOGS_URL = "http://localhost:{GATEWAY_PORT}/admin-panel/logs/{page_number}?pageSize={page_size}"
SAVE_DUMP_URL = "http://localhost:{GATEWAY_PORT}/admin-panel/save-dump"
LOAD_DUMP_URL = "http://localhost:{GATEWAY_PORT}/admin-panel/load-dump"

def fetch_user_profile():
    headers = {"Authorization": f"Bearer {st.session_state.get('jwt_token')}"}
    try:
        response = requests.get(PROFILE_URL.format(GATEWAY_PORT=GATEWAY_PORT), headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error fetching profile: {e}")
        return None

def download():
    st.warning("This website doesn't support downloading because this is personal project with no commercial potential")
    st.warning("Making downloading available will violate intellectual property rights")

def fetch_logs(page_number, page_size=10):
    headers = {"Authorization": f"Bearer {st.session_state.get('jwt_token')}"}
    try:
        response = requests.get(LOGS_URL.format(GATEWAY_PORT=GATEWAY_PORT, page_number=page_number, page_size=page_size), headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error fetching logs: {e}")
        return None

def extract_filename_from_header(header):
    if header and 'filename=' in header:
        fname = re.findall("filename=(.+)", header)
        if len(fname) > 0:
            return fname[0]
    return None

def save_dump():
    headers = {"Authorization": f"Bearer {st.session_state.get('jwt_token')}"}
    try:
        response = requests.get(SAVE_DUMP_URL.format(GATEWAY_PORT=GATEWAY_PORT), headers=headers)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        st.error(f"Error saving dump: {e}")
        return False
    
    content_disposition = response.headers.get('Content-Disposition')
    file_name = extract_filename_from_header(content_disposition)

    if not file_name:
        st.error("Failed to extract filename from the response header.")
        return
    
    with open(file_name, "wb") as file:
        for chunk in response.iter_content(chunk_size=8192):
            file.write(chunk)

    with open(file_name, "rb") as file:
        st.download_button(
            label="Download Database Dump",
            data=file,
            file_name=file_name,
            mime="application/sql"
        )
    return True

def load_dump(uploaded_file):
    headers = {"Authorization": f"Bearer {st.session_state.get('jwt_token')}"}

    if uploaded_file is not None:
        files = {
            'dump': (uploaded_file.name, uploaded_file.getvalue(), 'application/sql')
        }

        try:
            response = requests.post(LOAD_DUMP_URL.format(GATEWAY_PORT=GATEWAY_PORT), files=files, headers=headers)
            response.raise_for_status()
        except requests.exceptions.RequestException as e:
            st.error(f"Failed to upload database dump: {e}")
            return
        
        if "order" in st.session_state:
            st.session_state["order"] = fetch_orders()
        return True
    
    return False

def display_user_profile():
    user_data = fetch_user_profile()
    
    if user_data:
        user = user_data.get("user")
        if user:
            st.header(f"User Profile: {user['nickname']}")
            st.image(get_image_url(user["imageURL"]), width=150)
            st.write(f"**Email:** {user['email']}")
            st.write(f"**Role:** {'Admin' if user.get('isAdmin') else 'User'}")
            st.write(f"**Balance:** {user['balance']}")

            if user_data.get("purchasedAlbums"):
                st.subheader("My Albums")
                album_chunks = list(chunked(user_data["purchasedAlbums"], 5))
                for row in album_chunks:
                    album_columns = st.columns(len(row))
                    for album, col in zip(row, album_columns):
                        with col:
                            st.image(get_image_url(album["imageURL"]), caption=album["name"], use_container_width=True)
                            link = f"[View Album Profile](?entity=albums&id={album['id']})"
                            st.markdown(link, unsafe_allow_html=True)
                            st.button("Download", on_click=download, key=f"download_{user['id']}_{album['id']}")
            else:
                st.write("No albums found for this user.")

            if user['isAdmin'] == True:
                if st.button("Logs", key=f"logs_{user['id']}"):
                    if "logs_page" not in st.session_state:
                        st.session_state["logs_page"] = 1
                    logs = fetch_logs(st.session_state["logs_page"])
                    if not logs or not logs.get("logs"):
                        st.warning("No logs found!")
                        return

                    st.session_state["logs_count"] = logs["logsCount"]

                    for log in logs["logs"]:
                        st.write(f"{log['buyer']['id']}:{log['buyer']['nickname']} - {log['album']['id']}:{log['album']['name']} - {log['loggingTime']}")
                    
                    st.button("Prev", key=f"logs_{st.session_state.get('jwt_token')}_prev", on_click=decrement_page, disabled=st.session_state["logs_page"] == 1)
                    st.button("Next", key=f"logs_{st.session_state.get('jwt_token')}_next", on_click=increment_page, disabled=st.session_state["logs_page"] * 10 >= st.session_state["logs_count"])
                if st.button("Save Dump", key=f"save_dump_{user['id']}"):
                    save_dump()
                uploaded_file = st.file_uploader("Upload a .sql file", type=["sql"])
                load_dump(uploaded_file)
        else:
            st.error("Unable to retrieve user data.")
    else:
        st.error("Could not load profile data.")

def decrement_page():
    st.session_state["logs_page"] -= 1
    
def increment_page():
    st.session_state["logs_page"] += 1

display_user_profile()
