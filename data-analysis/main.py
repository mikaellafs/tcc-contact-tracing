import db
import graph
import processing
import operator

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

# Process data
df = get_from_csv("data/data-1669162316639.csv")
df = processing.parse_timestamp_types(df)

constant_df = processing.aggregate_constant_contacts(df)
metrics_df = processing.get_contact_metrics(constant_df)

# Draw graphs
metrics = metrics_df.reset_index()

G = graph.build_graph_from_df(metrics, "avg_distance")
graph.draw_with_weights(G, *graph.split_edge_sizes(G, operator.le, 200), filename='wgraph-distance-avg.png')


G = graph.build_graph_from_df(metrics, "max_duration")
graph.draw_with_weights(G, *graph.split_edge_sizes(G, operator.ge, 60*15), filename='wgraph-duration-max.png')

data_risk_columns = ["avg_distance", "max_duration"]
G = graph.build_graph_from_df(metrics, "avg_distance", data_risk_columns)

at_risk_nodes, not_risk_nodes = graph.split_contact_risks(G, *data_risk_columns)
graph.draw_with_weights(G, at_risk_nodes, not_risk_nodes, filename='wgraph-risk.png')
graph.draw_with_weights(G, at_risk_nodes, not_risk_nodes, filename='wgraph-risk-labeled.png', labels=True)

