package permissions

func shiftIndex[S ~[]E, E any](slice S, s int) S {
	return append(slice[:s], slice[s+1:]...)
}
