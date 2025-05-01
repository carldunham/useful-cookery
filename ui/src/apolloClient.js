import { ApolloClient, InMemoryCache, createHttpLink, from } from "@apollo/client";
import { onError } from "@apollo/client/link/error";

// Determine the GraphQL API URL based on the environment
const apiUrl = import.meta.env.PROD
  ? "/api/graphql" // Production URL (relative to the domain)
  : "http://localhost:8080/graphql"; // Development URL

console.log(`Using GraphQL API URL: ${apiUrl}`);

// Create an http link to the GraphQL API
const httpLink = createHttpLink({
  uri: apiUrl,
  credentials: "include", // Include cookies for authentication if needed
});

// Error handling link to log GraphQL and network errors
const errorLink = onError(({ graphQLErrors, networkError, operation }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path }) => {
      console.error(
        `[GraphQL error in operation ${operation.operationName}]: Message: ${message}, Location: ${locations}, Path: ${path}`
      );
    });
  }
  if (networkError) {
    console.error(`[Network error in operation ${operation.operationName}]: ${networkError}`);
  }
});

// Create the Apollo Client
const client = new ApolloClient({
  link: from([errorLink, httpLink]), // Chain the error link and http link
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
