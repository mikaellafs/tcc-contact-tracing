import networkx as nx
import matplotlib.pyplot as plt

def build_graph_from_df(data, isMulti=False):
    """ Build graph from dataframe (user1, user2, weight)"""
    G =  nx.Graph()

    for _, contact in data.iterrows():
        if not isMulti and (G.has_edge(contact["user1"], contact["user2"]) or G.has_edge(contact["user2"], contact["user1"])):
            continue

        G.add_edge(contact["user1"], contact["user2"], contact["weight"])
    
    return G

def draw(G, labels=False, filename='graph.png'):
    """ Draw a graph """
    # Compute position of nodes
    # pos = nx.spring_layout(G, k=0.5)
    pos = nx.kamada_kawai_layout(G)

    # Draw nodes and edges
    nx.draw_networkx_nodes(G, pos)
    nx.draw_networkx_edges(
        G, pos,
        # connectionstyle='arc3, rad=0.2'
    )

    plt.savefig(make_file_path(filename))
    plt.show()

def make_file_path(filename):
    return 'generated/' + filename
    
