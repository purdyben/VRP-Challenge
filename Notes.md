Vehicle Routing Problem 

type VRP interface{
	Drivers()
	Loads()
}

within the duration of one 12-hour shift. Your program does not have to assess problem feasibility.
●	No problem will contain more than 200 loads.


Unbounded number of drivers. (D)

The total cost of a solution is given by the formula:
total_cost = 500*number_of_drivers + total_number_of_driven_minutes 
