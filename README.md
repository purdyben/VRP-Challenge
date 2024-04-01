Benjamin Purdy Vorto VTP Challange

#### Problem:
The vehicle routing problem (VRP) which asks "What is the optimal set of routes for a fleet of vehicles to traverse in order to deliver to a given set of customers?" It generalises the travelling salesman problem (TSP). 

Determining the optimal solution to VRP is NP-hard. Felt this 

#### Instructions How to Run:

Please Run ```make``` to build the go binary stored in /bin/main 

```make eval``` will build and run evaluateShared.py and test data trainingProblems
 
Alternatively please run ```make evalb args="--problemDir YOUR_FOLDER"```

to build and run your evaluation  

#### Project Structure: 

#### Solution: 

#### Interesting Notes:

#### Failers:
As every driver starts and end at (0,0) we will in turn create a weighted tree structure and with a cyclic component. 


There is a path from every point to (0,0)

Every Node is 


Optimization using trainingProblems: 

Problem 1 line 9. 
9 (-41.48405901129298,-139.38690997500595) (-82.99128121032932,73.38972329128366)

0,0 -> pickup -> dropoff -> 0.0 = dis: 473.00277744002494 making this a solo trip, 