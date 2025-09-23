class CounselingCenter {
  final String id;
  final String name;
  final String address;
  final String phone;
  final bool onlineCounseling;
  final String createdAt;
  final String updatedAt;

  CounselingCenter({
    required this.id,
    required this.name,
    required this.address,
    required this.phone,
    required this.onlineCounseling,
    required this.createdAt,
    required this.updatedAt,
  });

  factory CounselingCenter.fromJson(Map<String, dynamic> json) {
    return CounselingCenter(
      id: json['id'] as String,
      name: json['name'] as String,
      address: json['address'] as String,
      phone: json['phone'] as String,
      onlineCounseling: json['online_counseling'] as bool,
      createdAt: json['created_at'] as String,
      updatedAt: json['updated_at'] as String,
    );
  }
}
