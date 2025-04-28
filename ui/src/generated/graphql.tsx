import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  DateTime: { input: string; output: string; }
};

export type AuthPayload = {
  __typename?: 'AuthPayload';
  token: Scalars['String']['output'];
  user: User;
};

export type Category = {
  __typename?: 'Category';
  description: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  recipes: Maybe<Array<Recipe>>;
};

export type DetailedIngredient = {
  __typename?: 'DetailedIngredient';
  id: Scalars['ID']['output'];
  isOptional: Maybe<Scalars['Boolean']['output']>;
  name: Scalars['String']['output'];
  preparation: Maybe<Scalars['String']['output']>;
  quantity: Maybe<Scalars['Float']['output']>;
  substitutes: Maybe<Array<Scalars['String']['output']>>;
  unit: Maybe<Scalars['String']['output']>;
};

export type Image = {
  __typename?: 'Image';
  alt: Maybe<Scalars['String']['output']>;
  height: Maybe<Scalars['Int']['output']>;
  id: Scalars['ID']['output'];
  url: Scalars['String']['output'];
  width: Maybe<Scalars['Int']['output']>;
};

export type ImageInput = {
  alt: InputMaybe<Scalars['String']['input']>;
  height: InputMaybe<Scalars['Int']['input']>;
  url: Scalars['String']['input'];
  width: InputMaybe<Scalars['Int']['input']>;
};

export type IngredientInput = {
  isOptional: InputMaybe<Scalars['Boolean']['input']>;
  name: Scalars['String']['input'];
  preparation: InputMaybe<Scalars['String']['input']>;
  quantity: InputMaybe<Scalars['Float']['input']>;
  substitutes: InputMaybe<Array<Scalars['String']['input']>>;
  unit: InputMaybe<Scalars['String']['input']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  addReview: Review;
  createCategory: Category;
  createRecipe: Recipe;
  deleteCategory: Scalars['Boolean']['output'];
  deleteRecipe: Scalars['Boolean']['output'];
  deleteReview: Scalars['Boolean']['output'];
  likeRecipe: Recipe;
  login: AuthPayload;
  register: AuthPayload;
  saveRecipe: User;
  unsaveRecipe: User;
  updateCategory: Category;
  updateRecipe: Recipe;
  updateReview: Review;
  updateUser: User;
};


export type MutationAddReviewArgs = {
  comment: InputMaybe<Scalars['String']['input']>;
  rating: Scalars['Int']['input'];
  recipeID: Scalars['ID']['input'];
};


export type MutationCreateCategoryArgs = {
  description: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
};


export type MutationCreateRecipeArgs = {
  input: RecipeInput;
};


export type MutationDeleteCategoryArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteRecipeArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteReviewArgs = {
  id: Scalars['ID']['input'];
};


export type MutationLikeRecipeArgs = {
  id: Scalars['ID']['input'];
};


export type MutationLoginArgs = {
  email: Scalars['String']['input'];
  password: Scalars['String']['input'];
};


export type MutationRegisterArgs = {
  input: RegisterInput;
};


export type MutationSaveRecipeArgs = {
  id: Scalars['ID']['input'];
};


export type MutationUnsaveRecipeArgs = {
  id: Scalars['ID']['input'];
};


export type MutationUpdateCategoryArgs = {
  description: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
  name: InputMaybe<Scalars['String']['input']>;
};


export type MutationUpdateRecipeArgs = {
  id: Scalars['ID']['input'];
  input: RecipeInput;
};


export type MutationUpdateReviewArgs = {
  comment: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
  rating: InputMaybe<Scalars['Int']['input']>;
};


export type MutationUpdateUserArgs = {
  input: UpdateUserInput;
};

export type NutritionInfo = {
  __typename?: 'NutritionInfo';
  calories: Maybe<Scalars['Int']['output']>;
  carbs: Maybe<Scalars['Float']['output']>;
  fat: Maybe<Scalars['Float']['output']>;
  fiber: Maybe<Scalars['Float']['output']>;
  protein: Maybe<Scalars['Float']['output']>;
  sodium: Maybe<Scalars['Int']['output']>;
  sugar: Maybe<Scalars['Float']['output']>;
};

export type NutritionInfoInput = {
  calories: InputMaybe<Scalars['Int']['input']>;
  carbs: InputMaybe<Scalars['Float']['input']>;
  fat: InputMaybe<Scalars['Float']['input']>;
  fiber: InputMaybe<Scalars['Float']['input']>;
  protein: InputMaybe<Scalars['Float']['input']>;
  sodium: InputMaybe<Scalars['Int']['input']>;
  sugar: InputMaybe<Scalars['Float']['input']>;
};

export enum OrderDirection {
  Asc = 'ASC',
  Desc = 'DESC'
}

export type PageInfo = {
  __typename?: 'PageInfo';
  endCursor: Maybe<Scalars['String']['output']>;
  hasNextPage: Scalars['Boolean']['output'];
  hasPreviousPage: Scalars['Boolean']['output'];
  startCursor: Maybe<Scalars['String']['output']>;
};

export type Query = {
  __typename?: 'Query';
  categories: Array<Category>;
  category: Maybe<Category>;
  findSubstitutes: Array<Scalars['String']['output']>;
  me: Maybe<User>;
  recipe: Maybe<Recipe>;
  recipeByOriginalId: Maybe<Recipe>;
  recipes: RecipeConnection;
  recommendRecipes: RecipeConnection;
  searchRecipes: RecipeConnection;
  user: Maybe<User>;
};


export type QueryCategoryArgs = {
  id: Scalars['ID']['input'];
};


export type QueryFindSubstitutesArgs = {
  ingredientName: Scalars['String']['input'];
};


export type QueryRecipeArgs = {
  id: Scalars['ID']['input'];
};


export type QueryRecipeByOriginalIdArgs = {
  originalId: Scalars['String']['input'];
};


export type QueryRecipesArgs = {
  after: InputMaybe<Scalars['String']['input']>;
  filter: InputMaybe<RecipeFilter>;
  first: InputMaybe<Scalars['Int']['input']>;
  offset: InputMaybe<Scalars['Int']['input']>;
  order: InputMaybe<RecipeOrder>;
};


export type QueryRecommendRecipesArgs = {
  after: InputMaybe<Scalars['String']['input']>;
  availableIngredients: InputMaybe<Array<Scalars['String']['input']>>;
  first: InputMaybe<Scalars['Int']['input']>;
  userID: InputMaybe<Scalars['ID']['input']>;
};


export type QuerySearchRecipesArgs = {
  after: InputMaybe<Scalars['String']['input']>;
  first: InputMaybe<Scalars['Int']['input']>;
  query: Scalars['String']['input'];
};


export type QueryUserArgs = {
  id: Scalars['ID']['input'];
};

export type Recipe = {
  __typename?: 'Recipe';
  author: Maybe<User>;
  averageRating: Maybe<Scalars['Float']['output']>;
  categories: Maybe<Array<Category>>;
  cookTime: Maybe<Scalars['Int']['output']>;
  createdAt: Maybe<Scalars['DateTime']['output']>;
  cuisine: Maybe<Scalars['String']['output']>;
  description: Maybe<Scalars['String']['output']>;
  difficulty: Maybe<SkillLevel>;
  id: Scalars['ID']['output'];
  images: Maybe<Array<Image>>;
  ingredients: Maybe<Array<DetailedIngredient>>;
  likes: Maybe<Scalars['Int']['output']>;
  nutritionInfo: Maybe<NutritionInfo>;
  originalId: Maybe<Scalars['String']['output']>;
  prepTime: Maybe<Scalars['Int']['output']>;
  reviews: Maybe<Array<Review>>;
  savedBy: Maybe<Array<User>>;
  servings: Maybe<Scalars['Int']['output']>;
  steps: Maybe<Array<Step>>;
  tags: Maybe<Array<Scalars['String']['output']>>;
  title: Scalars['String']['output'];
  updatedAt: Maybe<Scalars['DateTime']['output']>;
};

export type RecipeConnection = {
  __typename?: 'RecipeConnection';
  edges: Array<RecipeEdge>;
  pageInfo: PageInfo;
  totalCount: Scalars['Int']['output'];
};

export type RecipeEdge = {
  __typename?: 'RecipeEdge';
  cursor: Scalars['String']['output'];
  node: Recipe;
};

export type RecipeFilter = {
  authorID: InputMaybe<Scalars['ID']['input']>;
  categories: InputMaybe<Array<Scalars['ID']['input']>>;
  cuisine: InputMaybe<Scalars['String']['input']>;
  difficulty: InputMaybe<SkillLevel>;
  ingredients: InputMaybe<Array<Scalars['String']['input']>>;
  maxPrepTime: InputMaybe<Scalars['Int']['input']>;
  originalId: InputMaybe<Scalars['String']['input']>;
  search: InputMaybe<Scalars['String']['input']>;
};

export type RecipeInput = {
  categoryIDs: InputMaybe<Array<Scalars['ID']['input']>>;
  cookTime: InputMaybe<Scalars['Int']['input']>;
  cuisine: InputMaybe<Scalars['String']['input']>;
  description: Scalars['String']['input'];
  difficulty: InputMaybe<SkillLevel>;
  images: InputMaybe<Array<ImageInput>>;
  ingredients: Array<IngredientInput>;
  nutritionInfo: InputMaybe<NutritionInfoInput>;
  originalId: InputMaybe<Scalars['String']['input']>;
  prepTime: InputMaybe<Scalars['Int']['input']>;
  servings: InputMaybe<Scalars['Int']['input']>;
  steps: Array<StepInput>;
  tags: InputMaybe<Array<Scalars['String']['input']>>;
  title: Scalars['String']['input'];
};

export type RecipeOrder = {
  direction: OrderDirection;
  field: RecipeOrderField;
};

export enum RecipeOrderField {
  CreatedAt = 'CREATED_AT',
  Likes = 'LIKES',
  Rating = 'RATING',
  Title = 'TITLE'
}

export type RegisterInput = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  password: Scalars['String']['input'];
};

export type Review = {
  __typename?: 'Review';
  author: User;
  comment: Maybe<Scalars['String']['output']>;
  createdAt: Maybe<Scalars['DateTime']['output']>;
  id: Scalars['ID']['output'];
  rating: Scalars['Int']['output'];
  recipe: Recipe;
  updatedAt: Maybe<Scalars['DateTime']['output']>;
};

export enum Role {
  Admin = 'ADMIN',
  User = 'USER'
}

export enum SkillLevel {
  Advanced = 'ADVANCED',
  Beginner = 'BEGINNER',
  Intermediate = 'INTERMEDIATE'
}

export type Step = {
  __typename?: 'Step';
  description: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  image: Maybe<Image>;
  orderIndex: Maybe<Scalars['Int']['output']>;
  timeEstimate: Maybe<Scalars['Int']['output']>;
};

export type StepInput = {
  description: Scalars['String']['input'];
  image: InputMaybe<ImageInput>;
  orderIndex: Scalars['Int']['input'];
  timeEstimate: InputMaybe<Scalars['Int']['input']>;
};

export type Subscription = {
  __typename?: 'Subscription';
  newReview: Maybe<Review>;
  recipeLikes: Maybe<Scalars['Int']['output']>;
};


export type SubscriptionNewReviewArgs = {
  recipeID: Scalars['ID']['input'];
};


export type SubscriptionRecipeLikesArgs = {
  id: Scalars['ID']['input'];
};

export type UpdateUserInput = {
  email: InputMaybe<Scalars['String']['input']>;
  name: InputMaybe<Scalars['String']['input']>;
  password: InputMaybe<Scalars['String']['input']>;
  preferences: InputMaybe<UserPreferencesInput>;
};

export type User = {
  __typename?: 'User';
  createdAt: Maybe<Scalars['DateTime']['output']>;
  createdRecipes: Maybe<Array<Recipe>>;
  email: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  preferences: Maybe<UserPreferences>;
  reviews: Maybe<Array<Review>>;
  role: Role;
  savedRecipes: Maybe<Array<Recipe>>;
  updatedAt: Maybe<Scalars['DateTime']['output']>;
};

export type UserPreferences = {
  __typename?: 'UserPreferences';
  cuisinePreferences: Maybe<Array<Scalars['String']['output']>>;
  dietaryRestrictions: Maybe<Array<Scalars['String']['output']>>;
  dislikedIngredients: Maybe<Array<Scalars['String']['output']>>;
  favoriteIngredients: Maybe<Array<Scalars['String']['output']>>;
  skillLevel: Maybe<SkillLevel>;
};

export type UserPreferencesInput = {
  cuisinePreferences: InputMaybe<Array<Scalars['String']['input']>>;
  dietaryRestrictions: InputMaybe<Array<Scalars['String']['input']>>;
  dislikedIngredients: InputMaybe<Array<Scalars['String']['input']>>;
  favoriteIngredients: InputMaybe<Array<Scalars['String']['input']>>;
  skillLevel: InputMaybe<SkillLevel>;
};

export type SearchRecipesQueryVariables = Exact<{
  query: Scalars['String']['input'];
  first: InputMaybe<Scalars['Int']['input']>;
  after: InputMaybe<Scalars['String']['input']>;
}>;


export type SearchRecipesQuery = { __typename?: 'Query', searchRecipes: { __typename?: 'RecipeConnection', totalCount: number, edges: Array<{ __typename?: 'RecipeEdge', cursor: string, node: { __typename?: 'Recipe', id: string, title: string, description: string | null, prepTime: number | null, cookTime: number | null, difficulty: SkillLevel | null, averageRating: number | null, likes: number | null, images: Array<{ __typename?: 'Image', url: string, alt: string | null }> | null, author: { __typename?: 'User', name: string } | null, categories: Array<{ __typename?: 'Category', name: string }> | null } }>, pageInfo: { __typename?: 'PageInfo', hasNextPage: boolean, hasPreviousPage: boolean, startCursor: string | null, endCursor: string | null } } };

export type RecommendRecipesQueryVariables = Exact<{
  userID: InputMaybe<Scalars['ID']['input']>;
  availableIngredients: InputMaybe<Array<Scalars['String']['input']> | Scalars['String']['input']>;
  first: InputMaybe<Scalars['Int']['input']>;
  after: InputMaybe<Scalars['String']['input']>;
}>;


export type RecommendRecipesQuery = { __typename?: 'Query', recommendRecipes: { __typename?: 'RecipeConnection', totalCount: number, edges: Array<{ __typename?: 'RecipeEdge', cursor: string, node: { __typename?: 'Recipe', id: string, title: string, description: string | null, prepTime: number | null, cookTime: number | null, difficulty: SkillLevel | null, averageRating: number | null, likes: number | null, images: Array<{ __typename?: 'Image', url: string, alt: string | null }> | null, author: { __typename?: 'User', name: string } | null, categories: Array<{ __typename?: 'Category', name: string }> | null } }>, pageInfo: { __typename?: 'PageInfo', hasNextPage: boolean, hasPreviousPage: boolean, startCursor: string | null, endCursor: string | null } } };

export type RecipeCardFragmentFragment = { __typename?: 'Recipe', id: string, title: string, description: string | null, prepTime: number | null, cookTime: number | null, difficulty: SkillLevel | null, averageRating: number | null, likes: number | null, images: Array<{ __typename?: 'Image', url: string, alt: string | null }> | null, author: { __typename?: 'User', name: string } | null, categories: Array<{ __typename?: 'Category', name: string }> | null };

export type ExampleSearchRecipesQueryVariables = Exact<{
  query: Scalars['String']['input'];
  first: InputMaybe<Scalars['Int']['input']>;
  after: InputMaybe<Scalars['String']['input']>;
}>;


export type ExampleSearchRecipesQuery = { __typename?: 'Query', searchRecipes: { __typename?: 'RecipeConnection', totalCount: number, edges: Array<{ __typename?: 'RecipeEdge', cursor: string, node: { __typename?: 'Recipe', id: string, title: string, description: string | null, prepTime: number | null, cookTime: number | null, difficulty: SkillLevel | null, averageRating: number | null, likes: number | null, images: Array<{ __typename?: 'Image', url: string, alt: string | null }> | null, author: { __typename?: 'User', name: string } | null, categories: Array<{ __typename?: 'Category', name: string }> | null } }>, pageInfo: { __typename?: 'PageInfo', hasNextPage: boolean, hasPreviousPage: boolean, startCursor: string | null, endCursor: string | null } } };

export type ExampleRecommendRecipesQueryVariables = Exact<{
  userID: InputMaybe<Scalars['ID']['input']>;
  availableIngredients: InputMaybe<Array<Scalars['String']['input']> | Scalars['String']['input']>;
  first: InputMaybe<Scalars['Int']['input']>;
  after: InputMaybe<Scalars['String']['input']>;
}>;


export type ExampleRecommendRecipesQuery = { __typename?: 'Query', recommendRecipes: { __typename?: 'RecipeConnection', totalCount: number, edges: Array<{ __typename?: 'RecipeEdge', cursor: string, node: { __typename?: 'Recipe', id: string, title: string, description: string | null, prepTime: number | null, cookTime: number | null, difficulty: SkillLevel | null, averageRating: number | null, likes: number | null, images: Array<{ __typename?: 'Image', url: string, alt: string | null }> | null, author: { __typename?: 'User', name: string } | null, categories: Array<{ __typename?: 'Category', name: string }> | null } }>, pageInfo: { __typename?: 'PageInfo', hasNextPage: boolean, hasPreviousPage: boolean, startCursor: string | null, endCursor: string | null } } };

export const RecipeCardFragmentFragmentDoc = gql`
    fragment RecipeCardFragment on Recipe {
  id
  title
  description
  prepTime
  cookTime
  difficulty
  averageRating
  likes
  images {
    url
    alt
  }
  author {
    name
  }
  categories {
    name
  }
}
    `;
export const SearchRecipesDocument = gql`
    query SearchRecipes($query: String!, $first: Int, $after: String) {
  searchRecipes(query: $query, first: $first, after: $after) {
    edges {
      node {
        id
        title
        description
        prepTime
        cookTime
        difficulty
        averageRating
        likes
        images {
          url
          alt
        }
        author {
          name
        }
        categories {
          name
        }
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
    `;

/**
 * __useSearchRecipesQuery__
 *
 * To run a query within a React component, call `useSearchRecipesQuery` and pass it any options that fit your needs.
 * When your component renders, `useSearchRecipesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSearchRecipesQuery({
 *   variables: {
 *      query: // value for 'query'
 *      first: // value for 'first'
 *      after: // value for 'after'
 *   },
 * });
 */
export function useSearchRecipesQuery(baseOptions: Apollo.QueryHookOptions<SearchRecipesQuery, SearchRecipesQueryVariables> & ({ variables: SearchRecipesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SearchRecipesQuery, SearchRecipesQueryVariables>(SearchRecipesDocument, options);
      }
export function useSearchRecipesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SearchRecipesQuery, SearchRecipesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SearchRecipesQuery, SearchRecipesQueryVariables>(SearchRecipesDocument, options);
        }
export function useSearchRecipesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SearchRecipesQuery, SearchRecipesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SearchRecipesQuery, SearchRecipesQueryVariables>(SearchRecipesDocument, options);
        }
export type SearchRecipesQueryHookResult = ReturnType<typeof useSearchRecipesQuery>;
export type SearchRecipesLazyQueryHookResult = ReturnType<typeof useSearchRecipesLazyQuery>;
export type SearchRecipesSuspenseQueryHookResult = ReturnType<typeof useSearchRecipesSuspenseQuery>;
export type SearchRecipesQueryResult = Apollo.QueryResult<SearchRecipesQuery, SearchRecipesQueryVariables>;
export const RecommendRecipesDocument = gql`
    query RecommendRecipes($userID: ID, $availableIngredients: [String!], $first: Int, $after: String) {
  recommendRecipes(
    userID: $userID
    availableIngredients: $availableIngredients
    first: $first
    after: $after
  ) {
    edges {
      node {
        id
        title
        description
        prepTime
        cookTime
        difficulty
        averageRating
        likes
        images {
          url
          alt
        }
        author {
          name
        }
        categories {
          name
        }
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
    `;

/**
 * __useRecommendRecipesQuery__
 *
 * To run a query within a React component, call `useRecommendRecipesQuery` and pass it any options that fit your needs.
 * When your component renders, `useRecommendRecipesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useRecommendRecipesQuery({
 *   variables: {
 *      userID: // value for 'userID'
 *      availableIngredients: // value for 'availableIngredients'
 *      first: // value for 'first'
 *      after: // value for 'after'
 *   },
 * });
 */
export function useRecommendRecipesQuery(baseOptions?: Apollo.QueryHookOptions<RecommendRecipesQuery, RecommendRecipesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<RecommendRecipesQuery, RecommendRecipesQueryVariables>(RecommendRecipesDocument, options);
      }
export function useRecommendRecipesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<RecommendRecipesQuery, RecommendRecipesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<RecommendRecipesQuery, RecommendRecipesQueryVariables>(RecommendRecipesDocument, options);
        }
export function useRecommendRecipesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<RecommendRecipesQuery, RecommendRecipesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<RecommendRecipesQuery, RecommendRecipesQueryVariables>(RecommendRecipesDocument, options);
        }
export type RecommendRecipesQueryHookResult = ReturnType<typeof useRecommendRecipesQuery>;
export type RecommendRecipesLazyQueryHookResult = ReturnType<typeof useRecommendRecipesLazyQuery>;
export type RecommendRecipesSuspenseQueryHookResult = ReturnType<typeof useRecommendRecipesSuspenseQuery>;
export type RecommendRecipesQueryResult = Apollo.QueryResult<RecommendRecipesQuery, RecommendRecipesQueryVariables>;
export const ExampleSearchRecipesDocument = gql`
    query ExampleSearchRecipes($query: String!, $first: Int, $after: String) {
  searchRecipes(query: $query, first: $first, after: $after) {
    edges {
      node {
        ...RecipeCardFragment
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
    ${RecipeCardFragmentFragmentDoc}`;

/**
 * __useExampleSearchRecipesQuery__
 *
 * To run a query within a React component, call `useExampleSearchRecipesQuery` and pass it any options that fit your needs.
 * When your component renders, `useExampleSearchRecipesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useExampleSearchRecipesQuery({
 *   variables: {
 *      query: // value for 'query'
 *      first: // value for 'first'
 *      after: // value for 'after'
 *   },
 * });
 */
export function useExampleSearchRecipesQuery(baseOptions: Apollo.QueryHookOptions<ExampleSearchRecipesQuery, ExampleSearchRecipesQueryVariables> & ({ variables: ExampleSearchRecipesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ExampleSearchRecipesQuery, ExampleSearchRecipesQueryVariables>(ExampleSearchRecipesDocument, options);
      }
export function useExampleSearchRecipesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ExampleSearchRecipesQuery, ExampleSearchRecipesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ExampleSearchRecipesQuery, ExampleSearchRecipesQueryVariables>(ExampleSearchRecipesDocument, options);
        }
export function useExampleSearchRecipesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<ExampleSearchRecipesQuery, ExampleSearchRecipesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<ExampleSearchRecipesQuery, ExampleSearchRecipesQueryVariables>(ExampleSearchRecipesDocument, options);
        }
export type ExampleSearchRecipesQueryHookResult = ReturnType<typeof useExampleSearchRecipesQuery>;
export type ExampleSearchRecipesLazyQueryHookResult = ReturnType<typeof useExampleSearchRecipesLazyQuery>;
export type ExampleSearchRecipesSuspenseQueryHookResult = ReturnType<typeof useExampleSearchRecipesSuspenseQuery>;
export type ExampleSearchRecipesQueryResult = Apollo.QueryResult<ExampleSearchRecipesQuery, ExampleSearchRecipesQueryVariables>;
export const ExampleRecommendRecipesDocument = gql`
    query ExampleRecommendRecipes($userID: ID, $availableIngredients: [String!], $first: Int, $after: String) {
  recommendRecipes(
    userID: $userID
    availableIngredients: $availableIngredients
    first: $first
    after: $after
  ) {
    edges {
      node {
        ...RecipeCardFragment
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
    ${RecipeCardFragmentFragmentDoc}`;

/**
 * __useExampleRecommendRecipesQuery__
 *
 * To run a query within a React component, call `useExampleRecommendRecipesQuery` and pass it any options that fit your needs.
 * When your component renders, `useExampleRecommendRecipesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useExampleRecommendRecipesQuery({
 *   variables: {
 *      userID: // value for 'userID'
 *      availableIngredients: // value for 'availableIngredients'
 *      first: // value for 'first'
 *      after: // value for 'after'
 *   },
 * });
 */
export function useExampleRecommendRecipesQuery(baseOptions?: Apollo.QueryHookOptions<ExampleRecommendRecipesQuery, ExampleRecommendRecipesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ExampleRecommendRecipesQuery, ExampleRecommendRecipesQueryVariables>(ExampleRecommendRecipesDocument, options);
      }
export function useExampleRecommendRecipesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ExampleRecommendRecipesQuery, ExampleRecommendRecipesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ExampleRecommendRecipesQuery, ExampleRecommendRecipesQueryVariables>(ExampleRecommendRecipesDocument, options);
        }
export function useExampleRecommendRecipesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<ExampleRecommendRecipesQuery, ExampleRecommendRecipesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<ExampleRecommendRecipesQuery, ExampleRecommendRecipesQueryVariables>(ExampleRecommendRecipesDocument, options);
        }
export type ExampleRecommendRecipesQueryHookResult = ReturnType<typeof useExampleRecommendRecipesQuery>;
export type ExampleRecommendRecipesLazyQueryHookResult = ReturnType<typeof useExampleRecommendRecipesLazyQuery>;
export type ExampleRecommendRecipesSuspenseQueryHookResult = ReturnType<typeof useExampleRecommendRecipesSuspenseQuery>;
export type ExampleRecommendRecipesQueryResult = Apollo.QueryResult<ExampleRecommendRecipesQuery, ExampleRecommendRecipesQueryVariables>;