import db
import processing
import pandas as pd

def generate_from_raw():
    df = pd.read_csv("data/data-1669162316639.csv")
    df = processing.parse_timestamp_types(df)

    constant_df = processing.aggregate_constant_contacts(df)
    metrics_df = processing.get_contact_metrics(constant_df)

    return constant_df, metrics_df

def get_from_generated():
    constant_df = pd.read_csv("generated/csv/constant_contacts.csv").set_index([*processing.users_col, "contact"])
    metrics_df = pd.read_csv("generated/csv/contact_metrics.csv").set_index(processing.users_col)

    return constant_df, metrics_df

def get_infected_users(filename="infected.csv"):
    df = pd.read_csv("data/" + filename)

    return df["userid"].to_list()

def get_notified_users(filename="notified.csv"):
    df = pd.read_csv("data/" + filename)

    return df["userid"].to_list()
