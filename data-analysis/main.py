import graph
import operator
import moreplots
import datasource
import processing

import pandas as pd

########### Process data ###################

# constant_df, metrics_df = datasource.generate_from_raw()
constant_df, metrics_df = datasource.get_from_generated()

metrics = metrics_df.reset_index()

# # # Set uuids to int
# current_names = set(metrics["user1"].unique()).union(set(metrics["user2"].unique()))
# new_names = processing.generate_new_names(current_names)

# # Set names
# metrics = processing.rename_users_df(metrics, new_names)

########### Draw graphs ###################
data_risk_columns = ["total", "avg_distance", "max_duration"]

def draw_graphs():
    # Draw total contacts
    G = graph.build_graph_from_df(metrics, "avg_distance", data_risk_columns)

    pos = graph.draw(G, filename= "graph-total.png")

    # Draw graph with average distance risk
    graph.draw_with_weights(G, *graph.split_edge_sizes(G, operator.le, 200), pos=pos, filename='wgraph-distance-avg.png')

    # Draw graph with max duration risk
    G = graph.build_graph_from_df(metrics, "max_duration", data_risk_columns)
    graph.draw_with_weights(G, *graph.split_edge_sizes(G, operator.ge, 15), filename='wgraph-duration-max.png', pos=pos)

    # Draw graph considering average distance and max duration to evaluate risk
    at_risk_nodes, not_risk_nodes = graph.split_contact_risks(G, *data_risk_columns[1:])
    graph.draw_with_weights(G, at_risk_nodes, not_risk_nodes, filename='wgraph-risk.png', pos=pos)
    graph.draw_with_weights(G, at_risk_nodes, not_risk_nodes, filename='wgraph-risk-labeled.png', labels=True, pos=pos)

    # Draw graph graph with infected users
    infected = datasource.get_infected_users()
    graph.draw_with_weights(G, at_risk_nodes, not_risk_nodes, filename='wgraph-risk-infected.png', pos=pos, infected_nodes=infected)

    # Draw graph with infected users and notified
    notified = datasource.get_notified_users()
    graph.draw_with_weights(G, at_risk_nodes, not_risk_nodes, filename='wgraph-risk-infected-notified.png', pos=pos, infected_nodes=infected, notified_nodes=notified)

    return G

########### Plot differents metrics ###################

#### Get contact with max amount of registers
idx_contact_with_max_registers = metrics["total"].idxmax()
contact_with_max_registers = metrics.loc[idx_contact_with_max_registers]

def plot_cdf(G):
    ### Cdf of contact with max amount of registers
    cdf_base_filename_max_registers = "cdf-" + contact_with_max_registers["user1"] + "-" + contact_with_max_registers["user2"]
    print(cdf_base_filename_max_registers)

    # Distance
    moreplots.plot_cdf(constant_df.loc[(contact_with_max_registers["user1"], contact_with_max_registers["user2"])], "distance", filename= cdf_base_filename_max_registers + "-distance.png", title="CDF distancia contato com maior quantidade de registros", x="Distância (cm)")

    # Duration
    moreplots.plot_cdf(constant_df.loc[(contact_with_max_registers["user1"], contact_with_max_registers["user2"])], "duration", filename= cdf_base_filename_max_registers + "-duration.png", title="CDF duração contato com maior quantidade de registros", x="Duração (min)")

    # Node's degree
    degrees = pd.DataFrame(G.degree(), columns=['user', 'degree'])
    moreplots.plot_cdf(degrees, "degree", x="Grau", filename="cdf-all-degree.png")

    ### Cdf of all contacts
    # Distance
    moreplots.plot_cdf(constant_df, "distance", filename= "cdf-all-distance.png", title="CDF distancia todos os contatos", x="Distância (cm)")

    # Duration
    moreplots.plot_cdf(constant_df, "duration", filename= "cdf-all-duration.png", title="CDF duração todos os contatos", x="Duração (min)")

#### Scatter
def plot_scatter():
    moreplots.scatter_for_contact(constant_df.loc[(contact_with_max_registers["user1"], contact_with_max_registers["user2"])], filename="scatter-max-registers.png", title="Distribuição de contatos maior registro")
    moreplots.scatter_for_contact(constant_df, filename="scatter-all.png", title="")

def print_metrics(df, values):
    for value in values:
        print("Métricas de ", value)
        print("Média: ", df[value].mean())
        print("Desvio padrão: ", df[value].std())
        print("Min e máx: ", df[value].min(), " ", df[value].max())
        print("====================\n")

# print_metrics(constant_df, ["distance", "duration"])
# print_metrics(metrics_df, ["total"])
# G = draw_graphs()
# plot_cdf(G)
# plot_scatter()

df = pd.read_csv("data/data-1669162316639.csv")
df = processing.parse_timestamp_types(df)
moreplots.contacts_by_time(df)

# graph.make_csv_edge_gephi(metrics, 'avg_distance')