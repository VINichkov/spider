package entity

type CompareResult struct {
	Status string
	CompanyId
}

func NewCompareResult() *CompareResult{
	return &CompareResult{
		Status:    "NoError",
		CompanyId: CompanyId{},
	}
}

func (c *CompareResult)NoError()bool  {
	result := false
	if c.Status == "NoError"{
		result = true
	}
	return result
}

func (c *CompareResult)ToBlockList()  {
	c.Status = "BlockList"
}

func (c *CompareResult)Is_BlockList()bool  {
	result := false
	if c.Status == "BlockList"{
		result = true
	}
	return result
}

func (c *CompareResult)ToSameSource()  {
	c.Status = "SameSource"
}

func (c *CompareResult)Is_SameSource()bool  {
	result := false
	if c.Status == "SameSource"{
		result = true
	}
	return result
}

func (c *CompareResult)ToSameTitle()  {
	c.Status = "SameTitle"
}

func (c *CompareResult)Is_SameTitle()bool  {
	result := false
	if c.Status == "SameTitle"{
		result = true
	}
	return result
}

func (c *CompareResult)ToAnyError()  {
	c.Status = "AnyError"
}

func (c *CompareResult)Is_AnyError()bool  {
	result := false
	if c.Status == "AnyError"{
		result = true
	}
	return result
}