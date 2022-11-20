import db
import graph
import processing

import pandas as pd

def get_from_db():
    # connect to db
    conn = db.connect()
    cursor = conn.cursor()

    # get contacts
    df = db.select_contacts(cursor)

    # close connection
    cursor.close()
    conn.close()

    return df

def get_from_csv(filepath):
    return pd.read_csv(filepath)


df = get_from_csv("data/data-1668650971716.csv")
df = processing.parse_timestamp_types(df)
processing.aggregate_constant_contacts(df)
    

