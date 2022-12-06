from configparser import ConfigParser
import psycopg2
import pandas as pd

def config(filename='database.ini', section='postgresql-local'):
    """ Get config infos """
    parser = ConfigParser()
    parser.read(filename)

    # get section
    db = {}
    if parser.has_section(section):
        params = parser.items(section)
        for param in params:
            db[param[0]] = param[1]
        print(db)
    else:
        raise Exception('Section {0} not found in the {1} file'.format(section, filename))
    
    return db

def connect(section='postgresql-local'):
    """ Connect to the PostgreSQL database server """
    dbconfig = config(section=section)
    
    conn = None
    try:
        # connect to the PostgreSQL server
        print('Connecting to the PostgreSQL database...')
        conn = psycopg2.connect(**dbconfig)
    except (Exception, psycopg2.DatabaseError) as error:
        print(error)
    
    return conn

def select_contacts(cursor):
    query = ("""SELECT userId, anotherUser, firstcontacttimestamp,lastcontacttimestamp, distance 
                FROM contacts""")
    return pd.read_sql(query, cursor)

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