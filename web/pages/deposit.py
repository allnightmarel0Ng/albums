import streamlit as st
import requests

DEPOSIT_URL = "http://localhost:8080/deposit"

from main import get_authorization_header

def deposit_money(amount, headers):

    payload = {"money": amount}
    try:
        response = requests.post(DEPOSIT_URL, json=payload, headers=headers)
        response.raise_for_status()
        st.success(f"Successfully deposited ${amount}!")
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
                st.write(f"Deposit of ${amount} completed.")
            else:
                st.write("Failed to complete deposit.")
        else:
            st.error("Please enter a valid amount greater than zero.")

display_deposit_page()
