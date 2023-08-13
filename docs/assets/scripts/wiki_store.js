var search_store = {
    keyword                  : "",
    objects                  : objects,
    object_codes             : object_codes,
    relations                : relations,
    selected_categories      : [],
    selected_types           : [],
    selected_related_code    : '',
    filtered_results         : [],
    filtered_categories_count: [],
    showing_results          : [],
    unshowing_results        : [],
    count_per_page           : 30,
    is_loading               : false,
    loading_timer            : null,

    /**
     * init
     */

    init() {
        var url_params = new URLSearchParams(location.search);
        var q = url_params.get("q"); // keyword
        var c = url_params.get("c"); // categories
        var r = url_params.get("r"); // relate code
        var t = url_params.get("t"); // types

        if (q !== null && q !== '') {
            this.keyword = q.trim();
        }
        if (c !== null && c !== '') {
            this.selected_categories = c.trim().replace(/\s/g, "").split(",");
        }
        if (t !== null && t !== '') {
            this.selected_types = t.trim().replace(/\s/g, "").split(",");
        }
        if (r !== null && r !== '') {
            this.selected_related_code = r.trim().replace(/\s/g, "");
        }
        this.search();
    },

    /**
     * Reset Filter
     */

    resetFilter() {
        this.selected_categories = []
        this.updateParams()
        this.search()
    },

    /**
     * updateParams
     */

    resetRelatedCode() {
        this.selected_related_code = '';
        this.updateParams()
        this.search()
    },

    /**
     * updateParams
     */

    updateParams() {
        var searchParams = new URLSearchParams(window.location.search);

        if (this.selected_categories.length > 0) {
            searchParams.set("c", this.selected_categories.join(","));
        } else {
            searchParams.delete("c")
        }

        if (this.selected_types.length > 0) { // TODO: Showing on UI
            searchParams.set("t", this.selected_types.join(","));
        } else {
            searchParams.delete("t")
        }

        if (this.selected_related_code !== '') {
            searchParams.set("r", this.selected_related_code);
        } else {
            searchParams.delete("r")
        }

        if (this.keyword.trim() !== '') {
             searchParams.set("q", this.keyword.trim());
        } else {
            searchParams.delete("q")
        }
        window.history.replaceState(null, null, "?" + searchParams.toString());
    },

    /**
     * getGrandCodeIDs
     */

    getGrandCodeIDs(id) {
        if (id === -1) {
            return []
        }
        ids = [id]
        if (this.objects[id].g !== -1) { // g = Grand Code
            ids = [...ids, ...this.getGrandCodeIDs(this.objects[id].g)]
        }
        return ids
    },

    /**
     * searchDatabase
     */

    searchDatabase() {
        this.filtered_results = this.objects;

        if (this.selected_types.length > 0) {
            this.filtered_results = this.filtered_results.filter(obj => {
                // t = Type IDs
                return obj.t.map(i => relations.types[i]).some(v => this.selected_types.includes(v))
            })
        }

        if (this.selected_related_code !== '') {
            var code_id = object_codes.indexOf(this.selected_related_code) // if -1
            var grand_ids = this.getGrandCodeIDs(code_id)

            this.filtered_results = this.filtered_results.filter(obj => {
                // g = Grand Code ID, c = Code ID
                return (grand_ids.includes(obj.g) || grand_ids.includes(obj.c)) && obj.c !== code_id
            })
        }

        relations.categories.forEach((v, k) => {
            this.filtered_categories_count[k] = this.filtered_results.filter(v => v.r.includes(k)).length
        })

        if (this.selected_categories.length > 0) {
            this.filtered_results = this.filtered_results.filter(obj => {
                // r = Category Code IDs
                return obj.r.map(i => relations.categories[i]).some(v => this.selected_categories.includes(v))
            })
        }
    },

    /**
     * search
     */

    search() {
        this.is_loading = true;
        this.updateParams();

        this.unshowing_results = [];
        this.showing_results = [];

        if (this.loading_timer !== null) {
            clearTimeout(this.loading_timer);
        }
        this.loading_timer = setTimeout(() => {
            this.searchDatabase();

            if (this.keyword === '') {
                this.unshowing_results = this.filtered_results
            } else {
                var fuse = new Fuse(this.filtered_results, {
                    keys: ['n', 's'], // n = Name, s = Secondary Name
                })
                this.unshowing_results = fuse.search(this.keyword.trim()).map(v => v.item)
            }

            this.showMore(true);
            this.is_loading = false
        }, Math.random() * (300 - 200) + 300);
    },

    /**
     * showMore
     */

    showMore(ignore_loading) {
        if (ignore_loading === true) {
            this.showing_results = [...this.showing_results, ...this.unshowing_results.slice(0, this.count_per_page)]
            this.unshowing_results = this.unshowing_results.slice(this.count_per_page)
            return
        }
        this.is_loading = true;

        if (this.loading_timer !== null) {
            clearTimeout(this.loading_timer);
        }

        this.loading_timer = setTimeout(() => {
            this.showing_results = [...this.showing_results, ...this.unshowing_results.slice(0, this.count_per_page)]
            this.unshowing_results = this.unshowing_results.slice(this.count_per_page)
            this.is_loading = false
        }, Math.random() * (300 - 200) + 300);
    }
}