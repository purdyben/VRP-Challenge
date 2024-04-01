Benjamin Purdy Vorto VTP Challange

#### Problem:
The vehicle routing problem (VRP) asks "What is the optimal set of routes for a fleet of vehicles to traverse to deliver to a given set of customers?" It generalizes the traveling salesman problem (TSP). 

Determining the optimal solution to VRP is NP-hard. Felt this 

#### Instructions How to Run:

Build:
- Please Run ```make``` to build the go binary stored in /bin/main 

Training Data
- ```make eval``` will build and run evaluateShared.py and test data training problems
 
To run manually 
- build binary 
- ```python3 evaluateShared.py --cmd ./bin/main --problem YOUR_FOLDER ```

#### Project Structure: 

- cmd/main.go is the main file for the project 

#### Approach: 

My idea
- Cluster the given points based on approximation 
 	- Merge Clustering 
 	- Kmean Clustering
- Find Paths within Clusters 
- Repeat with different starting values to find the best approximation 
- Use goroutines to thread the computation

#### Training Data Testing 
Best Score 
- mean cost: 47759.118072847494
- mean run time: 279.8213050478981ms


