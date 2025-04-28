import { ApolloClient, InMemoryCache, createHttpLink } from "@apollo/client";

// Create an http link to the GraphQL API
const httpLink = createHttpLink({
  uri: "/api/graphql",
});

// Create the Apollo Client
const client = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: "network-only",
      errorPolicy: "all",
    },
    query: {
      fetchPolicy: "network-only",
      errorPolicy: "all",
    },
    mutate: {
      errorPolicy: "all",
    },
  },
});

export default client;
