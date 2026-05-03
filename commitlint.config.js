export default {
  extends: ["@commitlint/config-conventional"],
  rules: {
    "scope-enum": [
      2,
      "always",
      [
        "BR-1", "BR-2", "BR-3", "BR-4", "BR-5",
        "BR-6", "BR-7", "BR-8", "BR-9", "BR-10",
        "BR-11", "BR-12", "BR-13", "BR-14", "BR-15",
      ],
    ],
    "scope-case": [2, "always", "upper-case"],
  },
};
