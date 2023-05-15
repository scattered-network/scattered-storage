package helpers

// List
/* rbd --pool rbd list --format json
[
  "test-image1",
  "test-image2",
  "test-image3",
  "test-image4",
]
List is used for RBD images, RADOS pools, and other list worthy variables. */
type List []string
