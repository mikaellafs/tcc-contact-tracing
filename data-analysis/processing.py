import pandas as pd
import numpy as np
from sklearn import preprocessing

datetime_format='%Y-%m-%d %H:%M:%S.%f%z'
datetime_timezone='America/Sao_Paulo'

users_col = ['user1', 'user2']
timestamps_col = ['firstcontacttimestamp', 'lastcontacttimestamp']

max_total_seconds = 60*15 # 15 minutos

generated_base_filepath = 'generated'

def parse_timestamp_types(df):
    """ Parse datatypes in a dataframe of contacts 
    (..., firstcontacttimestamp,lastcontacttimestamp, ...)"""

    # Parse timestamps
    df['firstcontacttimestamp'] = pd.to_datetime(df['firstcontacttimestamp'], 
                                format=datetime_format, utc=True).dt.tz_convert(datetime_timezone)

    df['lastcontacttimestamp'] = pd.to_datetime(df['lastcontacttimestamp'], 
                                format=datetime_format, utc=True).dt.tz_convert(datetime_timezone)
    
    # Set column name
    df = df.rename(columns = {'userid':'user1', 'anotheruser': 'user2'})

    return df

def aggregate_constant_contacts(df, filename="constant_contacts.csv"):
    """ Aggregate contacts """

    # Label constant contacts
    d, label_col = label_constant_contacts(df)

    # Group constant contact to get metrics
    d = d.reset_index().groupby(label_col).agg(
        {
            "firstcontacttimestamp": 'min',
            "lastcontacttimestamp": 'max',
            "distance": 'mean'
        }
    )

    # Calculate duration of contacts
    d['duration'] = d.apply(lambda x : (x['lastcontacttimestamp'] - x['firstcontacttimestamp']).total_seconds(), axis=1)

    # Save to csv
    d.to_csv(make_file_path(filename))

    return d

def label_constant_contacts(df):
    """ Define constant contacts by labeling it with multiindex """
    # Sort values by contact timestamps
    df = df.sort_values(by=timestamps_col)

    # New df with constant contacts
    output, label_col = new_constant_contact_df()

    # label contacts
    for _, contact in df.iterrows():
        output = label_contact(output, contact)

    return output, label_col

def label_contact(df, contact):
    """ Label contact """
    idx = (contact['user1'], contact['user2'])

    if idx in df.index:
        df = check_constant_contact(df, idx, contact)
    else:
        df = new_constant_contact_row(df, idx, contact)
    
    return df

def new_constant_contact_df():
    """ Create an empty dataframe for constant contacts with multiindex (user1, user2, contact) and columns (timestamps, distance) """
    new_df_columns = [*users_col, 'contact', *timestamps_col, 'distance']
    new_idx_columns = [*users_col, 'contact']

    output = pd.DataFrame(columns= new_df_columns)
    output = output.set_index(new_idx_columns)

    return output, new_idx_columns

def new_constant_contact_row(df, users, contact, label=1):
    """ Add a new constant contact to dataframe """
    new_idx = (*users, label)

    row = pd.Series(contact[df.columns], name=new_idx)

    new_row = pd.DataFrame([row], columns=df.columns)
    df = pd.concat([df, new_row])
    
    return df.sort_index()

def check_constant_contact(df, users, contact):
    """ Check if contact is continuation of another one """
    # Get last user contacts
    last_contacts_idx = df.loc[users].index.get_level_values(0)[-1]
    last_contacts = df.loc[(*users, last_contacts_idx)]

    # Compute time diff
    last_timestamp = last_contacts['lastcontacttimestamp'].iloc[-1] if isinstance(last_contacts, pd.DataFrame) else last_contacts['lastcontacttimestamp']
    time_diff = contact['lastcontacttimestamp'] - last_timestamp

    # Check if constant contact
    is_constant = time_diff.total_seconds() <= max_total_seconds
    contact_id = last_contacts_idx if is_constant else last_contacts_idx +1

    df = new_constant_contact_row(df, users, contact, contact_id)

    return df

def get_contact_metrics(df, filename="contact_metrics.csv"):
    # Metrics about contact duration
    df = df.reset_index().groupby(users_col).agg(
        total=("contact", "count"),
        avg_duration=("duration", "mean"),
        std_duration=("duration", "std"),
        max_duration=("duration", "max"),
        min_duration=("duration", "min"),

        avg_distance=("distance", "mean"),
        std_distance=("distance", "std"),
        max_distance=("distance", "max"),
        min_distance=("distance", "min"),
    )

    # Save metrics to csv
    df.to_csv(make_file_path(filename))

    return df

def normalize(df):
    x = df.values #returns a numpy array
    min_max_scaler = preprocessing.MinMaxScaler()
    x_scaled = min_max_scaler.fit_transform(x.reshape(-1,1))
    df = pd.DataFrame(x_scaled)

    return df

def make_file_path(filename):
    return 'generated/csv/' + filename