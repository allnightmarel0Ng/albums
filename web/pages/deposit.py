import streamlit as st
import requests

from main import GATEWAY_PORT

DEPOSIT_URL = "http://localhost:{GATEWAY_PORT}/deposit"

from main import get_authorization_header

def deposit_money(amount, headers):
    payload = {"money": amount}
    try:
        response = requests.post(DEPOSIT_URL.format(GATEWAY_PORT=GATEWAY_PORT), json=payload, headers=headers)
        response.raise_for_status()
        return True
    except requests.exceptions.RequestException as e:
        st.error(f"Error depositing money: {e}")
        return False

def display_deposit_page():
    st.title("Deposit Money")

    headers = get_authorization_header()
    if not headers:
        st.error("You must be logged in to deposit money.")
        return False
    
    amount = st.number_input("Enter the amount to deposit:", min_value=1, step=1)

    if st.button("Deposit"):
        if amount > 0:
            if deposit_money(amount, headers):
                st.success(f"Deposit of ${amount} process has been started. Wait for completion in notifications!")
            else:
                st.write("Failed to start transaction!")
        else:
            st.error("Please enter a valid amount greater than zero.")

display_deposit_page()
