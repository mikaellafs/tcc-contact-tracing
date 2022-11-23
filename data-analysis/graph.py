import networkx as nx
import matplotlib.pyplot as plt

min_duration_risky=15*60 # 15 minutes
min_distance_risky=200 # 2 meters

# TODO: CONSIDER BOTH DIRECTIONS OR GET THE GREATEST ONE
def build_graph_from_df(data, weight_name, args, isMulti=False):
    """ Build graph from dataframe (user1, user2, weight)"""
    if isMulti:
        G = nx.MultiDiGraph()
    else:
        G = nx.Graph()

    for _, contact in data.iterrows():
        if not isMulti and (G.has_edge(contact["user1"], contact["user2"]) or G.has_edge(contact["user2"], contact["user1"])):
            continue
        
        G.add_edge(contact["user1"], contact["user2"], weight=contact[weight_name], data=contact[args].to_dict())
    
    return G

def draw(G, labels=False, filename='graph5.png'):
    """ Draw a graph """
    # Compute position of nodes
    pos = nx.spring_layout(G, k=10, seed=5)

    # Draw nodes and edges
    if labels:
        nx.draw_networkx_labels(G, pos)

    nx.draw_networkx_nodes(G, pos)
    nx.draw_networkx_edges(
        G, pos,
        connectionstyle='arc3, rad=0.2'
    )

    plt.savefig(make_file_path(filename))
    plt.show()

def draw_with_weights(G, elarge, esmall, labels=False, filename='wgraph.png'):
    pos = nx.spring_layout(G, k=10, seed=5)

    # Draw nodes
    if labels:
        nx.draw_networkx_labels(G, pos)

    nx.draw_networkx_nodes(G, pos)

    # Draw edges with weights
    nx.draw_networkx_edges(G, pos, edgelist=elarge, width=2)
    nx.draw_networkx_edges(
    G, pos, edgelist=esmall, width=2, alpha=0.5, edge_color="b", style="dashed"
    )

    plt.savefig(make_file_path(filename))
    plt.show()

def split_edge_sizes(G, operator, criteria):
    elarge = [(u, v) for (u, v, d) in G.edges(data=True) if operator(d["weight"], criteria)]
    esmall = [(u, v) for (u, v, d) in G.edges(data=True) if not operator(d["weight"], criteria)]

    return elarge, esmall

def split_contact_risks(G, distance, duration):
    elarge = [(u, v) for (u, v, d) in G.edges(data=True) if is_at_risk(d['data'][distance], d['data'][duration])]
    esmall = [(u, v) for (u, v, d) in G.edges(data=True) if not is_at_risk(d['data'][distance], d['data'][duration])]

    return elarge, esmall

def is_at_risk(distance, duration):
    return distance <=min_distance_risky and duration >= min_duration_risky

def make_file_path(filename):
    return 'generated/graph/' + filename