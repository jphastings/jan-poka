//=====================================
// This is public Domain Code
// Contributed by: William A Adams
// 11 May 2011
//=====================================
joinfactor = 0.125;

gSteps = 4;
gHeight = 4;


//=======================================
// Functions
// These are the 4 blending functions for a cubic bezier curve
//=======================================

/*
	Bernstein Basis Functions

	For Bezier curves, these functions give the weights per control point.
*/
function BEZ03(u) = pow((1-u), 3);
function BEZ13(u) = 3*u*(pow((1-u),2));
function BEZ23(u) = 3*(pow(u,2))*(1-u);
function BEZ33(u) = pow(u,3);

// Calculate a singe point along a cubic bezier curve
// Given a set of 4 control points, and a parameter 0 <= 'u' <= 1
// These functions will return the exact point on the curve
function PointOnBezCubic2D(p0, p1, p2, p3, u) = [
	BEZ03(u)*p0[0]+BEZ13(u)*p1[0]+BEZ23(u)*p2[0]+BEZ33(u)*p3[0],
	BEZ03(u)*p0[1]+BEZ13(u)*p1[1]+BEZ23(u)*p2[1]+BEZ33(u)*p3[1]];

