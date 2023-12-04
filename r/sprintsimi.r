# #Uncomment for memory profiling
 gc()
 Rprof("Rprof.out", memory.profiling=TRUE)

library(dplyr)
library(purrr)


# read
read_csv_to_df <- function(csv_path) {
  df <- read.csv(csv_path)
  return(df)
}

# add sum
add_sum_column <- function(df) {
  df %>%
    mutate(Sum = rowSums(select(., c("Jeremy", "Angelica", "Mike"))))
}

# percentile
monte_carlo_percentile <- function(column, percentile, iterations) {
  samples <- replicate(iterations, sample(column, size = length(column), replace = TRUE))
  quantile(as.numeric(samples), probs = percentile)
}

# iterations to goal
monte_carlo_iterations_to_goal <- function(column, goal, iterations) {
  points_needed <- replicate(iterations, {
    total <- 0
    sprints <- 0
    while(total < goal) {
      total <- total + sample(column, size = 1)
      sprints <- sprints + 1
    }
    sprints
  })
  mean(points_needed)
}

# main
csv_path <- "./data/sprints.csv"
df <- read_csv_to_df(csv_path)

# add sum column
df_with_sum <- add_sum_column(df)

# print the table
print(df_with_sum)

# monte carlo sims
jer50 <- monte_carlo_percentile(df_with_sum$Jeremy, 0.5, 1000)
ang50 <- monte_carlo_percentile(df_with_sum$Angelica, 0.5, 1000)
mik50 <- monte_carlo_percentile(df_with_sum$Mike, 0.5, 1000)
sum50 <- monte_carlo_percentile(df_with_sum$Sum, 0.5, 1000)
sum100pts <- monte_carlo_iterations_to_goal(df_with_sum$Sum, 100, 1000)

cat("Jeremy's 50th percentile:", jer50, "\n")
cat("Angelica's 50th percentile:", ang50, "\n")
cat("Mike's 50th percentile:", mik50, "\n")
cat("Team's 50th percentile:", sum50, "\n")
cat("Average number of sprints to complete 100 pts:", sum100pts, "\n")

# # Uncomment for memory profiling
 Rprof(NULL)
 max(summaryRprof("Rprof.out", memory="both")$by.total$mem.total)

