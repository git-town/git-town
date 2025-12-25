// https://docs.rs/mdbook/latest/mdbook/book/struct.Book.html
export type Book = {
  sections: BookItem[]
}

// https://docs.rs/mdbook/latest/mdbook/book/enum.BookItem.html
export type BookItem =
  | { Chapter: Chapter }
  | "Separator"
  | { PartTitle: string }

// https://docs.rs/mdbook/latest/mdbook/book/struct.Chapter.html
export type Chapter = {
  name: string
  content: string
  number?: SectionNumber
  sub_items: BookItem[]
  path?: string
  source_path?: string
  parent_names: string[]
}

// https://docs.rs/mdbook/latest/mdbook/book/struct.SectionNumber.html
export type SectionNumber = number[]
